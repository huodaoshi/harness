package configpaths

import (
	"path/filepath"
	"runtime"
)

// BackendRoot returns the backend module root (directory containing go.mod).
func BackendRoot() string {
	_, file, _, _ := runtime.Caller(0)
	// internal/configpaths -> backend/
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func SafetyRules() string {
	return filepath.Join(BackendRoot(), "config", "safety_rules_v1.yaml")
}

func CrisisTemplatesZH() string {
	return filepath.Join(BackendRoot(), "config", "crisis_templates", "zh-CN.json")
}

func BoundaryTemplatesZH() string {
	return filepath.Join(BackendRoot(), "config", "boundary_templates", "zh-CN.json")
}

// WebRoot returns the harness web static assets directory (sibling of backend/).
func WebRoot() string {
	return filepath.Clean(filepath.Join(BackendRoot(), "..", "web"))
}
