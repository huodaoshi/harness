package nextchat_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/api/nextchat"
)

func TestProxyHandler_ForwardsRequestBody(t *testing.T) {
	var gotBody []byte
	upstream := http.NewServeMux()
	upstream.HandleFunc("/api/v3/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upstream body: %v", err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	})
	upSrv := &http.Server{Handler: upstream}
	upLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go upSrv.Serve(upLn)
	t.Cleanup(func() { _ = upSrv.Close() })

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	h := server.New(server.WithListener(ln))
	settings := nextchat.Settings{
		ARKAPIKey:  "test-key",
		ARKModel:   "ep-rewritten",
		ARKBaseURL: "http://" + upLn.Addr().String(),
	}
	proxy := &nextchat.ProxyHandler{Settings: settings, Client: http.DefaultClient}
	h.POST("/api/bytedance/*path", proxy.Handle)
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(30 * time.Millisecond)

	payload := map[string]any{
		"model":  "Harness",
		"stream": true,
		"messages": []map[string]string{
			{"role": "user", "content": "hello"},
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		"http://"+ln.Addr().String()+"/api/bytedance/api/v3/chat/completions",
		"application/json",
		bytes.NewReader(raw),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d body %s", resp.StatusCode, b)
	}

	if len(gotBody) == 0 {
		t.Fatal("upstream received empty body")
	}
	var upstreamPayload map[string]any
	if err := json.Unmarshal(gotBody, &upstreamPayload); err != nil {
		t.Fatalf("upstream body not json: %q err=%v", gotBody, err)
	}
	if upstreamPayload["model"] != "ep-rewritten" {
		t.Fatalf("model = %v, want ep-rewritten", upstreamPayload["model"])
	}
}
