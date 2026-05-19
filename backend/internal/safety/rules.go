package safety

import (
	"os"

	"github.com/huodaoshi/harness/backend/internal/configpaths"
	"gopkg.in/yaml.v3"
)

// Rule is one L1 safety rule group.
type Rule struct {
	ID       string   `yaml:"id"`
	Patterns []string `yaml:"patterns"`
}

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

func defaultRulesPath() string {
	return configpaths.SafetyRules()
}

// LoadRules reads safety rules from a YAML file.
func LoadRules(path string) ([]Rule, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f rulesFile
	if err := yaml.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	return f.Rules, nil
}
