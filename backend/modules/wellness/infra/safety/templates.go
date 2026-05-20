package safety

import (
	"encoding/json"
	"os"

	"github.com/huodaoshi/harness/backend/modules/wellness/infra/configpaths"
)

// CrisisPayload is returned to clients on the crisis SSE branch.
type CrisisPayload struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

type crisisTemplate struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

// TemplateStore loads fixed crisis copy by template id.
type TemplateStore struct {
	byID map[string]crisisTemplate
}

func defaultTemplatesPath() string {
	return configpaths.CrisisTemplatesZH()
}

// NewTemplateStore loads bundled zh-CN crisis templates.
func NewTemplateStore() (*TemplateStore, error) {
	return LoadTemplates(defaultTemplatesPath())
}

// LoadTemplates reads crisis templates JSON.
func LoadTemplates(path string) (*TemplateStore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]crisisTemplate
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &TemplateStore{byID: raw}, nil
}

// Render returns crisis SSE payload for a gate result.
func (s *TemplateStore) Render(r Result) (CrisisPayload, bool) {
	if !r.IsCrisis() {
		return CrisisPayload{}, false
	}
	t, ok := s.byID[r.TemplateID]
	if !ok {
		return CrisisPayload{}, false
	}
	return CrisisPayload{TemplateID: t.TemplateID, Body: t.Body}, true
}
