package conf_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/huodaoshi/harness/backend/conf"
)

func TestLoad_localDefaults(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	backendRoot := filepath.Clean(filepath.Join(wd, "..", ".."))
	t.Chdir(backendRoot)

	t.Setenv("APP_ENV", "local")
	t.Setenv("USE_MEMORY_STORE", "true")

	cfg, err := conf.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.App.Port != 8080 {
		t.Fatalf("port=%d want 8080", cfg.App.Port)
	}
	if cfg.MongoDB.Database != "family_wellness" {
		t.Fatalf("db=%q", cfg.MongoDB.Database)
	}
	if !cfg.Redis.Required {
		t.Fatal("base config should require redis for auth")
	}
}
