package knowledgeindexing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/schema"

	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
)

const (
	maxChunkRunes     = 1200
	chunkOverlapRunes = 100
)

// Parsers holds format-specific parsers (P1: text/markdown/json only).
type Parsers struct{}

// NewParsers creates parsers for supported formats.
func NewParsers(_ context.Context) (*Parsers, error) {
	return &Parsers{}, nil
}

func parseSource(_ context.Context, _ *Parsers, src ResolvedSource, job *domain.IngestJob) ([]*schema.Document, error) {
	meta := map[string]any{
		"space_id":    fmt.Sprintf("%d", job.SpaceID),
		"job_id":      job.JobID,
		"doc_type":    job.DocType,
		"source_type": job.SourceType,
		"doc_key":     job.DocKey,
		"source_url":  firstNonEmpty(job.SourceURL, src.FinalURL),
	}
	meta["source_format"] = string(src.Format)

	text := strings.TrimSpace(string(src.Body))
	if text == "" {
		return nil, nil
	}

	switch src.Format {
	case SourceFormatMarkdown, SourceFormatText:
		return []*schema.Document{{Content: text, MetaData: meta}}, nil
	case SourceFormatJSON:
		return parseJSONDocument(text, meta)
	default:
		return []*schema.Document{{Content: text, MetaData: meta}}, nil
	}
}

func parseJSONDocument(text string, meta map[string]any) ([]*schema.Document, error) {
	var v any
	if err := json.Unmarshal([]byte(text), &v); err != nil {
		return nil, fmt.Errorf("knowledgeindexing: json parse: %w", err)
	}
	var lines []string
	flattenJSON("$", v, &lines, 0, 8, 500)
	if len(lines) == 0 {
		return nil, nil
	}
	content := strings.Join(lines, "\n")
	return []*schema.Document{{Content: content, MetaData: meta}}, nil
}

func flattenJSON(path string, v any, lines *[]string, depth, maxDepth, maxNodes int) {
	if depth > maxDepth || len(*lines) >= maxNodes {
		return
	}
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			flattenJSON(path+"."+k, child, lines, depth+1, maxDepth, maxNodes)
		}
	case []any:
		for i, child := range val {
			flattenJSON(fmt.Sprintf("%s[%d]", path, i), child, lines, depth+1, maxDepth, maxNodes)
		}
	default:
		if val != nil {
			*lines = append(*lines, fmt.Sprintf("%s: %v", path, val))
		}
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// unused but keeps processor compatible if extended
var _ io.Reader
