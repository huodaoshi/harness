package apitest

import (
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/huodaoshi/harness/backend/api"
	authpkg "github.com/huodaoshi/harness/backend/modules/auth/pkg"
)

// Fixed guest UUID for API integration tests.
const testGuestAnonID = "11111111-1111-4111-8111-111111111111"

func testGuestUserID() string {
	return "anon:" + testGuestAnonID
}

func testJWTSigner(t *testing.T) authpkg.JWTSigner {
	t.Helper()
	return authpkg.NewHS256Signer("harness-test-jwt-secret", 3600)
}

func testGuestMw(t *testing.T) app.HandlerFunc {
	t.Helper()
	return api.JWTOrGuestMiddleware(testJWTSigner(t))
}

func withGuest(req *http.Request) {
	req.Header.Set("X-Anon-ID", testGuestAnonID)
}
