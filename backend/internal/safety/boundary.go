package safety

import (
	"encoding/json"
	"os"

	"github.com/huodaoshi/harness/backend/internal/configpaths"
)

// MedicalPayload is returned on the medical_boundary SSE branch.
type MedicalPayload struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

// BlockPayload is returned on the block SSE branch (PRD error shape).
type BlockPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type medicalTemplate struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

type blockTemplate struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BoundaryStore loads fixed medical and block copy.
type BoundaryStore struct {
	medical medicalTemplate
	block   blockTemplate
}

func defaultBoundaryPath() string {
	return configpaths.BoundaryTemplatesZH()
}

// NewBoundaryStore loads bundled zh-CN boundary templates.
func NewBoundaryStore() (*BoundaryStore, error) {
	return LoadBoundaryStore(defaultBoundaryPath())
}

// LoadBoundaryStore reads boundary templates JSON.
func LoadBoundaryStore(path string) (*BoundaryStore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Medical medicalTemplate `json:"medical_boundary"`
		Block   blockTemplate   `json:"block"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &BoundaryStore{medical: raw.Medical, block: raw.Block}, nil
}

// RenderMedical returns medical SSE payload.
func (s *BoundaryStore) RenderMedical(r Result) (MedicalPayload, bool) {
	if !r.IsMedical() {
		return MedicalPayload{}, false
	}
	return MedicalPayload{
		TemplateID: s.medical.TemplateID,
		Body:       s.medical.Body,
	}, s.medical.Body != ""
}

// RenderBlock returns block refusal payload.
func (s *BoundaryStore) RenderBlock(r Result) (BlockPayload, bool) {
	if !r.IsBlock() {
		return BlockPayload{}, false
	}
	return BlockPayload{
		Code:    s.block.Code,
		Message: s.block.Message,
	}, s.block.Message != ""
}
