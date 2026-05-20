package authpkg_test

import (
	"testing"

	authpkg "github.com/huodaoshi/harness/backend/modules/auth/pkg"
)

func TestGenerateJWTSecret(t *testing.T) {
	a, err := authpkg.GenerateJWTSecret()
	if err != nil {
		t.Fatal(err)
	}
	b, err := authpkg.GenerateJWTSecret()
	if err != nil {
		t.Fatal(err)
	}
	if len(a) == 0 || len(b) == 0 {
		t.Fatal("empty secret")
	}
	if a == b {
		t.Fatal("expected two different secrets")
	}
}
