package knowledgeindexing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"

	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
)

const maxStoredErrLen = 2000

// Processor consumes knowledge-ingest messages: fetch ? parse ? pipeline ? update status.
type Processor struct {
	rdb        redis.UniversalClient
	jobs       domain.IngestJobRepo
	spaces     domain.SpaceRepo
	fetch      FetchOptions
	objStorage infra.ObjectStorage
	pipeline   compose.Runnable[[]*schema.Document, []string]
	parsers    *Parsers
}

// NewProcessor constructs a Processor with injected dependencies.
func NewProcessor(
	rdb redis.UniversalClient,
	jobs domain.IngestJobRepo,
	spaces domain.SpaceRepo,
	fetch FetchOptions,
	objStorage infra.ObjectStorage,
	pipeline compose.Runnable[[]*schema.Document, []string],
	parsers *Parsers,
) *Processor {
	return &Processor{
		rdb:        rdb,
		jobs:       jobs,
		spaces:     spaces,
		fetch:      fetch,
		objStorage: objStorage,
		pipeline:   pipeline,
		parsers:    parsers,
	}
}

// Process handles one ingest MQ payload.
func (p *Processor) Process(ctx context.Context, payload []byte) error {
	var job domain.IngestJob
	if err := json.Unmarshal(payload, &job); err != nil {
		return fmt.Errorf("knowledgeindexing: unmarshal job: %w", err)
	}
	slog.InfoContext(ctx, "knowledgeindexing: process start",
		"job_id", job.JobID, "space_id", job.SpaceID,
		"source_type", job.SourceType, "source_url", job.SourceURL)

	if job.DocKey == "" {
		key, err := domain.BuildIngestDocKey(job.SpaceID, job.SourceType, job.SourceURL, job.Content)
		if err != nil {
			_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateStoredErr(err.Error()))
			return fmt.Errorf("knowledgeindexing: build doc_key: %w", err)
		}
		job.DocKey = key
	}

	_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusProcessing, "")

	raw, contentType, finalURL, err := p.resolveRaw(ctx, &job)
	if err != nil {
		_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateStoredErr(err.Error()))
		return fmt.Errorf("knowledgeindexing: resolve: %w", err)
	}

	src := resolveSource(raw, job.SourceURL, contentType, finalURL)
	docs, err := parseSource(ctx, p.parsers, src, &job)
	if err != nil {
		_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateStoredErr(err.Error()))
		return fmt.Errorf("knowledgeindexing: parse: %w", err)
	}
	if len(docs) == 0 {
		msg := "no indexable text extracted"
		_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, msg)
		return fmt.Errorf("knowledgeindexing: %s", msg)
	}

	slog.InfoContext(ctx, "knowledgeindexing: pipeline invoking",
		"job_id", job.JobID, "input_docs", len(docs))
	chunkIDs, err := p.pipeline.Invoke(ctx, docs)
	if err != nil {
		_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateStoredErr(err.Error()))
		return fmt.Errorf("knowledgeindexing: pipeline: %w", err)
	}

	oldCount, err := updateDocChunkMapping(ctx, p.rdb, job.SpaceID, job.DocKey, chunkIDs)
	if err != nil {
		_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateStoredErr(err.Error()))
		return fmt.Errorf("knowledgeindexing: update doc chunk mapping: %w", err)
	}

	delta := int64(len(chunkIDs)) - int64(oldCount)
	if p.spaces != nil && delta != 0 {
		if adjErr := p.spaces.AdjustChunkCount(ctx, job.SpaceID, delta); adjErr != nil {
			slog.WarnContext(ctx, "knowledgeindexing: adjust chunk count",
				"space_id", job.SpaceID, "delta", delta, "error", adjErr)
		}
	}
	if cntErr := p.jobs.UpdateChunkCount(ctx, job.JobID, len(chunkIDs)); cntErr != nil {
		slog.WarnContext(ctx, "knowledgeindexing: update chunk count", "job_id", job.JobID, "error", cntErr)
	}

	_ = p.jobs.UpdateStatus(ctx, job.JobID, domain.IngestStatusDone, "")
	slog.InfoContext(ctx, "knowledgeindexing: process done",
		"job_id", job.JobID, "space_id", job.SpaceID,
		"indexed_chunks", len(chunkIDs), "delta", delta)
	return nil
}

func (p *Processor) resolveRaw(ctx context.Context, job *domain.IngestJob) (body []byte, contentType, finalURL string, err error) {
	switch job.SourceType {
	case 1:
		return []byte(job.Content), "", "", nil
	case 2:
		if job.SourceURL == "" {
			return nil, "", "", fmt.Errorf("empty source_url")
		}
		return fetchSource(ctx, job.SourceURL, p.fetch)
	case 3:
		if job.SourceURL == "" {
			return nil, "", "", fmt.Errorf("empty source_url for source_type=3")
		}
		data, dlErr := p.objStorage.Download(ctx, job.SourceURL)
		return data, "", job.SourceURL, dlErr
	default:
		return nil, "", "", fmt.Errorf("unsupported source_type %d", job.SourceType)
	}
}

func truncateStoredErr(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxStoredErrLen {
		return s
	}
	return s[:maxStoredErrLen] + "?"
}
