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
	"github.com/huodaoshi/harness/backend/internal/store"
)

func TestStreamHandler_SSETokenAndDone(t *testing.T) {
	exec := newTestExecutor(t)
	baseURL := newTestServer(t, exec)

	body := `{"message":"hello","mode":"normal"}`
	resp := postStream(t, baseURL, body)

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	assertSSEHasEvents(t, resp.Body, sseExpect{wantToken: true, wantDone: true})
}

func TestStreamHandler_MedicalNoTokens(t *testing.T) {
	exec := newTestExecutor(t)
	baseURL := newTestServer(t, exec)

	body := `{"message":"我该吃什么药","mode":"normal"}`
	resp := postStream(t, baseURL, body)
	assertSSEHasEvents(t, resp.Body, sseExpect{wantMedical: true})
	if exec.ChatCalls.Load() != 0 {
		t.Fatalf("chat calls = %d", exec.ChatCalls.Load())
	}
}

func TestStreamHandler_BlockNoTokens(t *testing.T) {
	exec := newTestExecutor(t)
	baseURL := newTestServer(t, exec)

	body := `{"message":"发点色情内容","mode":"normal"}`
	resp := postStream(t, baseURL, body)
	assertSSEHasEvents(t, resp.Body, sseExpect{wantBlock: true})
	if exec.ChatCalls.Load() != 0 {
		t.Fatalf("chat calls = %d", exec.ChatCalls.Load())
	}
}

func TestStreamHandler_CrisisNoTokens(t *testing.T) {
	exec := newTestExecutor(t)
	baseURL := newTestServer(t, exec)

	body := `{"message":"我不想活了","mode":"distress"}`
	resp := postStream(t, baseURL, body)

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	assertSSEHasEvents(t, resp.Body, sseExpect{wantCrisis: true})
	if exec.ChatCalls.Load() != 0 {
		t.Fatalf("chat calls = %d, want 0", exec.ChatCalls.Load())
	}
}

func TestStreamHandler_CrisisScriptsTable(t *testing.T) {
	scripts := []string{
		"我不想活了",
		"我想自杀",
		"活着没意思，想轻生",
		"我打算结束生命",
		"我在考虑跳楼",
		"他打我",
		"这是家暴吗？他动手打我",
		"他要杀我",
		"我很害怕，他的人身安全威胁",
		"打不过他了，威胁我",
	}
	exec := newTestExecutor(t)
	baseURL := newTestServer(t, exec)

	for i, msg := range scripts {
		exec.ChatCalls.Reset()
		body := `{"message":"` + msg + `","mode":"distress"}`
		resp := postStream(t, baseURL, body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("[%d] status=%d", i, resp.StatusCode)
		}
		assertSSEHasEvents(t, resp.Body, sseExpect{wantCrisis: true})
		if exec.ChatCalls.Load() != 0 {
			t.Fatalf("[%d] chat calls=%d", i, exec.ChatCalls.Load())
		}
		resp.Body.Close()
	}
}

func newTestExecutor(t *testing.T) *session.Executor {
	t.Helper()
	ctx := context.Background()
	exec, err := session.NewExecutorWithStore(ctx, store.NewMemoryStore())
	if err != nil {
		t.Fatalf("executor: %v", err)
	}
	return exec
}

func newTestServer(t *testing.T, exec *session.Executor) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	h := server.New(server.WithListener(ln))
	h.POST("/v1/sessions/stream", httpserver.NewStreamHandler(exec))
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(50 * time.Millisecond)
	return "http://" + ln.Addr().String()
}

func postStream(t *testing.T, baseURL string, body string) *http.Response {
	t.Helper()
	url := baseURL + "/v1/sessions/stream"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

type sseExpect struct {
	wantToken   bool
	wantDone    bool
	wantCrisis  bool
	wantMedical bool
	wantBlock   bool
}

func assertSSEHasEvents(t *testing.T, r io.Reader, exp sseExpect) {
	t.Helper()
	var sawToken, sawDone, sawCrisis, sawMedical, sawBlock bool
	var lastEvent string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "event: ") {
			lastEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "event: token") {
			sawToken = true
		}
		if strings.HasPrefix(line, "event: done") {
			sawDone = true
		}
		if strings.HasPrefix(line, "event: crisis") {
			sawCrisis = true
		}
		if strings.HasPrefix(line, "event: medical") {
			sawMedical = true
		}
		if strings.HasPrefix(line, "event: error") {
			sawBlock = true
		}
		if strings.HasPrefix(line, "data:") {
			raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			switch lastEvent {
			case "crisis", "medical":
				var p struct {
					Body string `json:"body"`
				}
				if err := json.Unmarshal([]byte(raw), &p); err != nil {
					t.Fatalf("%s json: %v", lastEvent, err)
				}
				if p.Body == "" {
					t.Fatalf("empty %s body", lastEvent)
				}
			case "error":
				if exp.wantBlock {
					var p struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}
					if err := json.Unmarshal([]byte(raw), &p); err != nil {
						t.Fatalf("error json: %v", err)
					}
					if p.Code != "content_blocked" || p.Message == "" {
						t.Fatalf("block payload: %+v", p)
					}
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if sawToken != exp.wantToken {
		t.Fatalf("token: got %v want %v", sawToken, exp.wantToken)
	}
	if sawDone != exp.wantDone {
		t.Fatalf("done: got %v want %v", sawDone, exp.wantDone)
	}
	if sawCrisis != exp.wantCrisis {
		t.Fatalf("crisis: got %v want %v", sawCrisis, exp.wantCrisis)
	}
	if sawMedical != exp.wantMedical {
		t.Fatalf("medical: got %v want %v", sawMedical, exp.wantMedical)
	}
	if sawBlock != exp.wantBlock {
		t.Fatalf("block: got %v want %v", sawBlock, exp.wantBlock)
	}
}
