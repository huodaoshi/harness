package httpserver_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/internal/httpserver"
	"github.com/huodaoshi/harness/backend/internal/session"
)

func TestStreamHandler_SSETokenAndDone(t *testing.T) {
	ctx := context.Background()
	runnable, err := session.NewStreamGraph(ctx)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	h := server.New(server.WithListener(ln))
	h.POST("/v1/sessions/stream", httpserver.NewStreamHandler(runnable))

	go h.Spin()
	defer func() {
		_ = h.Close()
	}()
	time.Sleep(50 * time.Millisecond)

	body := `{"message":"hello","mode":"normal"}`
	url := "http://" + ln.Addr().String() + "/v1/sessions/stream"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("content-type=%s", resp.Header.Get("Content-Type"))
	}

	var sawToken, sawDone bool
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "event: token") {
			sawToken = true
		}
		if strings.HasPrefix(line, "event: done") {
			sawDone = true
		}
		if strings.HasPrefix(line, "data:") && sawToken && !sawDone {
			raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			var p struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal([]byte(raw), &p); err != nil {
				t.Fatalf("token json: %v", err)
			}
			if p.Text == "" {
				t.Fatal("empty token")
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if !sawToken {
		t.Fatal("missing token event")
	}
	if !sawDone {
		t.Fatal("missing done event")
	}
}
