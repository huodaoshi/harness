package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/huodaoshi/harness/backend/infra"
	knowledgeapp "github.com/huodaoshi/harness/backend/modules/knowledge/application"
	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
)

type submitIngestRequest struct {
	SourceType int    `json:"source_type"`
	SourceURL  string `json:"source_url"`
	Content    string `json:"content"`
	DocType    int    `json:"doc_type"`
}

type submitIngestResponse struct {
	JobID   string `json:"job_id"`
	SpaceID int64  `json:"space_id"`
	Status  int    `json:"status"`
}

type ingestJobResponse struct {
	JobID      string    `json:"job_id"`
	SpaceID    int64     `json:"space_id"`
	SourceType int       `json:"source_type"`
	SourceURL  string    `json:"source_url,omitempty"`
	Content    string    `json:"content,omitempty"`
	DocType    int       `json:"doc_type"`
	DocKey     string    `json:"doc_key,omitempty"`
	Status     int       `json:"status"`
	Error      string    `json:"error,omitempty"`
	ChunkCount int       `json:"chunk_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type searchHitResponse struct {
	ChunkID string  `json:"chunk_id"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
	DocType int     `json:"doc_type"`
}

// AdminKnowledgeHandler serves /v1/admin knowledge routes.
type AdminKnowledgeHandler struct {
	ingest    knowledgeapp.IngestService
	knowledge knowledgeapp.KnowledgeService
}

func NewAdminKnowledgeHandler(ingest knowledgeapp.IngestService, knowledge knowledgeapp.KnowledgeService) *AdminKnowledgeHandler {
	return &AdminKnowledgeHandler{ingest: ingest, knowledge: knowledge}
}

func RegisterAdminKnowledgeRoutes(r *server.Hertz, h *AdminKnowledgeHandler, jwtMw, adminMw app.HandlerFunc) {
	g := r.Group("/v1/admin")
	g.Use(jwtMw, adminMw)
	g.POST("/spaces/:space_id/ingest", h.handleSubmitIngest)
	g.GET("/ingest/:job_id", h.handleGetIngestJob)
	g.GET("/spaces/:space_id/search", h.handleSearch)
}

func (h *AdminKnowledgeHandler) handleSubmitIngest(ctx context.Context, c *app.RequestContext) {
	spaceID, ok := parseAdminSpaceID(ctx, c)
	if !ok {
		return
	}
	var req submitIngestRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}
	req.SourceURL = strings.TrimSpace(req.SourceURL)
	if !validKnowledgeDocType(req.DocType) {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}
	switch req.SourceType {
	case 1:
		if req.Content == "" || req.SourceURL != "" {
			apierror.Render(ctx, c, apierror.ErrBadRequest)
			return
		}
	case 2:
		if req.SourceURL == "" || req.Content != "" {
			apierror.Render(ctx, c, apierror.ErrBadRequest)
			return
		}
	default:
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

	job, err := h.ingest.SubmitIngestJob(ctx, spaceID, req.SourceType, req.SourceURL, req.Content, req.DocType, "")
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}
	c.JSON(consts.StatusCreated, submitIngestResponse{
		JobID:   job.JobID,
		SpaceID: job.SpaceID,
		Status:  job.Status,
	})
}

func (h *AdminKnowledgeHandler) handleGetIngestJob(ctx context.Context, c *app.RequestContext) {
	jobID := c.Param("job_id")
	if jobID == "" {
		apierror.Render(ctx, c, apierror.ErrNotFound)
		return
	}
	job, err := h.ingest.GetIngestJob(ctx, jobID)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, toIngestJobResp(job))
}

func (h *AdminKnowledgeHandler) handleSearch(ctx context.Context, c *app.RequestContext) {
	spaceID, ok := parseAdminSpaceID(ctx, c)
	if !ok {
		return
	}
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}
	topK := 5
	if raw := c.Query("top_k"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 20 {
			topK = n
		}
	}

	// Use fake/hash embedder path via KnowledgeService — caller must pass embedding.
	// For admin search we embed query with the same fake embedder used at index time.
	emb := infra.FakeEmbedVector(query)
	chunks, err := h.knowledge.RetrieveChunks(ctx, &spaceID, query, topK, emb)
	if err != nil {
		apierror.Render(ctx, c, apierror.ErrInternal)
		return
	}
	out := make([]searchHitResponse, 0, len(chunks))
	for _, ch := range chunks {
		out = append(out, searchHitResponse{
			ChunkID: ch.ChunkID,
			Content: ch.Content,
			Score:   ch.Score,
			DocType: ch.DocType,
		})
	}
	c.JSON(consts.StatusOK, map[string]any{"items": out})
}

func parseAdminSpaceID(ctx context.Context, c *app.RequestContext) (int64, bool) {
	raw := c.Param("space_id")
	var id int64
	if _, err := fmt.Sscanf(raw, "%d", &id); err != nil || id <= 0 {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return 0, false
	}
	return id, true
}

func validKnowledgeDocType(v int) bool {
	return v >= 1 && v <= 3
}

func toIngestJobResp(job *domain.IngestJob) ingestJobResponse {
	return ingestJobResponse{
		JobID:      job.JobID,
		SpaceID:    job.SpaceID,
		SourceType: job.SourceType,
		SourceURL:  job.SourceURL,
		Content:    job.Content,
		DocType:    job.DocType,
		DocKey:     job.DocKey,
		Status:     job.Status,
		Error:      job.Error,
		ChunkCount: job.ChunkCount,
		CreatedAt:  job.CreatedAt,
		UpdatedAt:  job.UpdatedAt,
	}
}
