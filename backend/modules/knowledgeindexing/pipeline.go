package knowledgeindexing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	redisindexer "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	redisCli "github.com/redis/go-redis/v9"

	knowledgeinf "github.com/huodaoshi/harness/backend/modules/knowledge/infra"
)

// BuildRedisIndexer initializes the eino-ext Redis indexer (keys: space_{id}:chunk:...).
func BuildRedisIndexer(
	ctx context.Context,
	rdbClient *redisCli.Client,
	embedder embedding.Embedder,
) (*knowledgeinf.Indexer, error) {
	cfg := &knowledgeinf.IndexerConfig{
		Client:    rdbClient,
		KeyPrefix: "",
		BatchSize: 24,
		Embedding: embedder,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*knowledgeinf.Hashes, error) {
			metaBytes, err := json.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("knowledgeindexing: indexer: marshal metadata: %w", err)
			}
			sourceType := metaInt(doc.MetaData, "source_type")
			docType := metaInt(doc.MetaData, "doc_type")
			sourceURL, _ := doc.MetaData["source_url"].(string)
			return &knowledgeinf.Hashes{
				Key: doc.ID,
				Field2Value: map[string]knowledgeinf.FieldValue{
					"content": {
						Value:    doc.Content,
						EmbedKey: "embedding",
					},
					"metadata": {Value: string(metaBytes)},
					"source_type": {Value: strconv.Itoa(sourceType)},
					"doc_type":    {Value: strconv.Itoa(docType)},
					"source_url":  {Value: sourceURL},
					"is_active":   {Value: "active"},
				},
			}, nil
		},
	}
	idx, err := knowledgeinf.NewIndexer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("knowledgeindexing: build redis indexer: %w", err)
	}
	return idx, nil
}

func metaInt(meta map[string]any, key string) int {
	if meta == nil {
		return 0
	}
	switch v := meta[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		return 0
	}
}

// BuildIndexingPipeline: resolve → parse → chunk → assignID → index (no LLM score/tag).
func BuildIndexingPipeline(
	ctx context.Context,
	indexer *knowledgeinf.Indexer,
) (compose.Runnable[[]*schema.Document, []string], error) {
	g := compose.NewGraph[[]*schema.Document, []string]()

	if err := g.AddLambdaNode("chunk", compose.InvokableLambda(
		func(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
			var out []*schema.Document
			for _, doc := range docs {
				format, _ := doc.MetaData["source_format"].(string)
				var chunks []*schema.Document
				if format == string(SourceFormatMarkdown) {
					chunks = splitMarkdownByHeaders(doc)
				} else {
					chunks = splitByRunes(doc, 1200, 100)
				}
				out = append(out, chunks...)
			}
			return out, nil
		},
	)); err != nil {
		return nil, fmt.Errorf("knowledgeindexing: pipeline: add chunk: %w", err)
	}

	if err := g.AddLambdaNode("assignID", compose.InvokableLambda(
		func(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
			for i, doc := range docs {
				spaceIDStr, _ := doc.MetaData["space_id"].(string)
				jobID, _ := doc.MetaData["job_id"].(string)
				spaceID, err := strconv.ParseInt(spaceIDStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("knowledgeindexing: assignID: parse space_id %q: %w", spaceIDStr, err)
				}
				doc.ID = knowledgeinf.BuildKnowledgeChunkKey(spaceID, jobID, i)
			}
			return docs, nil
		},
	)); err != nil {
		return nil, fmt.Errorf("knowledgeindexing: pipeline: add assignID: %w", err)
	}

	if err := g.AddIndexerNode("index", indexer); err != nil {
		return nil, fmt.Errorf("knowledgeindexing: pipeline: add index: %w", err)
	}

	for _, e := range []struct{ from, to string }{
		{compose.START, "chunk"},
		{"chunk", "assignID"},
		{"assignID", "index"},
		{"index", compose.END},
	} {
		if err := g.AddEdge(e.from, e.to); err != nil {
			return nil, fmt.Errorf("knowledgeindexing: pipeline: edge %s→%s: %w", e.from, e.to, err)
		}
	}

	runnable, err := g.Compile(ctx, compose.WithGraphName("KnowledgeIndexingPipelineHarness"))
	if err != nil {
		return nil, fmt.Errorf("knowledgeindexing: pipeline: compile: %w", err)
	}
	return runnable, nil
}

// updateDocChunkMapping updates live chunk set and deletes stale keys (best-effort).
func updateDocChunkMapping(
	ctx context.Context,
	rdb redisCli.UniversalClient,
	spaceID int64,
	docKey string,
	newChunkIDs []string,
) (oldCount int, err error) {
	liveKey := knowledgeinf.BuildDocSetKey(spaceID, docKey)

	oldKeys, err := rdb.SMembers(ctx, liveKey).Result()
	if err != nil {
		return 0, fmt.Errorf("knowledgeindexing: updateDocChunkMapping: smembers: %w", err)
	}
	oldCount = len(oldKeys)

	pipe := rdb.TxPipeline()
	pipe.Del(ctx, liveKey)
	if len(newChunkIDs) > 0 {
		members := make([]interface{}, len(newChunkIDs))
		for i, id := range newChunkIDs {
			members[i] = id
		}
		pipe.SAdd(ctx, liveKey, members...)
	}
	if _, execErr := pipe.Exec(ctx); execErr != nil {
		return 0, fmt.Errorf("knowledgeindexing: updateDocChunkMapping: exec: %w", execErr)
	}

	newSet := make(map[string]struct{}, len(newChunkIDs))
	for _, id := range newChunkIDs {
		newSet[id] = struct{}{}
	}
	var toDelete []string
	for _, k := range oldKeys {
		if _, exists := newSet[k]; !exists {
			toDelete = append(toDelete, k)
		}
	}
	if len(toDelete) > 0 {
		if delErr := knowledgeinf.DeleteChunkKeys(ctx, rdb, toDelete); delErr != nil {
			slog.WarnContext(ctx, "knowledgeindexing: delete old chunks", "space_id", spaceID, "error", delErr)
		}
	}
	return oldCount, nil
}

// splitMarkdownByHeaders splits on # headers (same semantics as one-eino).
func splitMarkdownByHeaders(doc *schema.Document) []*schema.Document {
	lines := strings.Split(doc.Content, "\n")
	var sections []string
	var cur strings.Builder
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			if cur.Len() > 0 {
				s := strings.TrimSpace(cur.String())
				if s != "" {
					sections = append(sections, s)
				}
				cur.Reset()
			}
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	if cur.Len() > 0 {
		s := strings.TrimSpace(cur.String())
		if s != "" {
			sections = append(sections, s)
		}
	}
	if len(sections) == 0 {
		sections = []string{strings.TrimSpace(doc.Content)}
	}
	return buildChildDocs(doc, sections)
}

func splitByRunes(doc *schema.Document, chunkSize, overlap int) []*schema.Document {
	text := strings.TrimSpace(doc.Content)
	if text == "" {
		return nil
	}
	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []*schema.Document{cloneDocWithContent(doc, text)}
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 4
	}
	step := chunkSize - overlap
	if step <= 0 {
		step = chunkSize
	}
	var chunks []string
	for start := 0; start < len(runes); start += step {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		if end == len(runes) {
			break
		}
	}
	return buildChildDocs(doc, chunks)
}

func buildChildDocs(parent *schema.Document, contents []string) []*schema.Document {
	out := make([]*schema.Document, 0, len(contents))
	for _, content := range contents {
		content = strings.TrimSpace(content)
		if content == "" {
			continue
		}
		out = append(out, &schema.Document{Content: content, MetaData: inheritMeta(parent.MetaData)})
	}
	return out
}

func cloneDocWithContent(parent *schema.Document, content string) *schema.Document {
	return &schema.Document{Content: content, MetaData: inheritMeta(parent.MetaData)}
}

func inheritMeta(parent map[string]any) map[string]any {
	if len(parent) == 0 {
		return make(map[string]any)
	}
	m := make(map[string]any, len(parent))
	for k, v := range parent {
		m[k] = v
	}
	return m
}

// silence unused import if redisindexer types change
var _ = redisindexer.IndexerConfig{}
