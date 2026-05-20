package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/infra/chatmodel"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

func TestSession_EndAndInject(t *testing.T) {
	st := store.NewMemoryStore()
	ctx := context.Background()
	exec, err := application.NewExecutorWithGateway(ctx, st, chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}
	baseURL := newProfileAndStreamServer(t, st, exec)
	uid := testGuestUserID()
	marker := "SUMMARY_LINE_ALPHA"

	putProfile(t, baseURL, uid, `{"self":{"note":""},"people":[],"current_issue":"`+marker+`"}`)

	body1 := `{"message":"第一轮","mode":"normal"}`
	resp1 := postStream(t, baseURL, body1)
	sid := parseSessionIDFromSSE(t, resp1.Body)
	resp1.Body.Close()

	endBody := `{"session_id":"` + sid + `"}`
	endResp := postJSON(t, baseURL+"/v1/sessions/end", endBody)
	if endResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(endResp.Body)
		t.Fatalf("end status=%d %s", endResp.StatusCode, b)
	}
	var endOut struct {
		Summary3 []string `json:"summary3"`
	}
	_ = json.NewDecoder(endResp.Body).Decode(&endOut)
	endResp.Body.Close()
	if len(endOut.Summary3) != 3 {
		t.Fatalf("summary3: %+v", endOut.Summary3)
	}

	body2 := `{"message":"新一场","mode":"normal"}`
	resp2 := postStream(t, baseURL, body2)
	full := readTokenText(t, resp2.Body)
	resp2.Body.Close()
	if !strings.Contains(full, endOut.Summary3[0]) && !strings.Contains(full, "[上次会话摘要]") {
		t.Fatalf("missing summary in stream: %q", full)
	}

	sess := getSessionHTTP(t, baseURL, sid, uid)
	if len(sess.Messages) < 2 {
		t.Fatalf("messages=%d", len(sess.Messages))
	}
}

func postJSON(t *testing.T, url, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	withGuest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func getSessionHTTP(t *testing.T, baseURL, sessionID, userID string) domain.StoredSession {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/sessions/"+sessionID, nil)
	if err != nil {
		t.Fatal(err)
	}
	withGuest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var s domain.StoredSession
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

func parseSessionIDFromSSE(t *testing.T, r io.Reader) string {
	t.Helper()
	full := readSSEEvents(t, r)
	for _, block := range full {
		if strings.Contains(block, "event: done") {
			for _, line := range strings.Split(block, "\n") {
				if strings.HasPrefix(line, "data:") {
					raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
					var p struct {
						SessionID string `json:"session_id"`
					}
					_ = json.Unmarshal([]byte(raw), &p)
					if p.SessionID != "" {
						return p.SessionID
					}
				}
			}
		}
	}
	t.Fatal("no session_id in SSE")
	return ""
}

func readTokenText(t *testing.T, r io.Reader) string {
	t.Helper()
	var parts []string
	for _, block := range readSSEEvents(t, r) {
		if !strings.Contains(block, "event: token") {
			continue
		}
		for _, line := range strings.Split(block, "\n") {
			if strings.HasPrefix(line, "data:") {
				raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				var p struct {
					Text string `json:"text"`
				}
				_ = json.Unmarshal([]byte(raw), &p)
				parts = append(parts, p.Text)
			}
		}
	}
	return strings.Join(parts, "")
}

func readSSEEvents(t *testing.T, r io.Reader) []string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(string(b), "\n\n")
}
