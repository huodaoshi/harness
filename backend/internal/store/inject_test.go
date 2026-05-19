package store_test

import (
	"testing"

	"github.com/huodaoshi/harness/backend/internal/store"
)

func TestBuildInjectBlock_Empty(t *testing.T) {
	if got := store.BuildInjectBlock(nil, nil); got != "" {
		t.Fatalf("got %q", got)
	}
}
