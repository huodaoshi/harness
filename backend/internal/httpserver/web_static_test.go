package httpserver_test

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/internal/configpaths"
	"github.com/huodaoshi/harness/backend/internal/httpserver"
)

func TestWebStatic_JSAndCSS(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	h := server.New(server.WithListener(ln))
	httpserver.RegisterWebStatic(h, configpaths.WebRoot())
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(50 * time.Millisecond)

	base := "http://" + ln.Addr().String()
	for _, path := range []string{
		"/js/app.js",
		"/js/sse.js",
		"/css/app.css",
		"/manifest.webmanifest",
	} {
		resp, err := http.Get(base + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s: status=%d body=%q", path, resp.StatusCode, body)
		}
	}

	resp, err := http.Get(base + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("/ status=%d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "disclaimer-check") {
		t.Fatal("index missing disclaimer checkbox")
	}
}
