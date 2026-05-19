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

	"github.com/huodaoshi/harness/backend/internal/chatmodel"
	"github.com/huodaoshi/harness/backend/internal/httpserver"
	"github.com/huodaoshi/harness/backend/internal/session"
	"github.com/huodaoshi/harness/backend/internal/store"
)

func TestProfile_GET_EmptyThenPUT_ReadBack(t *testing.T) {
	st := store.NewMemoryStore()
	baseURL := newProfileTestServer(t, st)
	uid := "user-profile-1"

	get := getProfile(t, baseURL, uid)
	if get.UserID != uid || get.CurrentIssue != "" {
		t.Fatalf("empty get: %+v", get)
	}
	if get.People == nil {
		t.Fatal("people should be [] not null")
	}

	putBody := `{"self":{"note":"容易焦虑"},"people":[{"label":"阿明","relation":"伴侣","note":""}],"current_issue":"反复冷战"}`
	put := putProfile(t, baseURL, uid, putBody)
	if put.CurrentIssue != "反复冷战" || put.Self.Note != "容易焦虑" {
		t.Fatalf("put: %+v", put)
	}

	got := getProfile(t, baseURL, uid)
	if got.CurrentIssue != put.CurrentIssue || len(got.People) != 1 {
		t.Fatalf("read back: %+v", got)
	}
}

func TestProfile_PUT_ThenStream_InjectsCurrentIssue(t *testing.T) {
	st := store.NewMemoryStore()
	ctx := context.Background()
	exec, err := session.NewExecutorWithGateway(ctx, st, chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}
	baseURL := newProfileAndStreamServer(t, st, exec)
	uid := "user-inject-6"
	marker := "INTEGRATION_ISSUE_XYZ"

	putProfile(t, baseURL, uid, `{"self":{"note":""},"people":[],"current_issue":"`+marker+`"}`)

	body := `{"user_id":"` + uid + `","message":"你好","mode":"normal"}`
	resp := postStream(t, baseURL, body)

	var tokens []string
	var lastEvent string
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "event: ") {
			lastEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data:") && lastEvent == "token" {
			raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			var p struct {
				Text string `json:"text"`
			}
			_ = json.Unmarshal([]byte(raw), &p)
			tokens = append(tokens, p.Text)
		}
	}
	full := strings.Join(tokens, "")
	if !strings.Contains(full, marker) {
		t.Fatalf("stream text missing %q: %q", marker, full)
	}
}

func TestProfile_RequiresUserID(t *testing.T) {
	st := store.NewMemoryStore()
	baseURL := newProfileTestServer(t, st)
	resp, err := http.Get(baseURL + "/v1/profile")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}

func newProfileTestServer(t *testing.T, st store.Store) string {
	t.Helper()
	return listenProfileRoutes(t, st, nil)
}

func newProfileAndStreamServer(t *testing.T, st store.Store, exec *session.Executor) string {
	t.Helper()
	return listenProfileRoutes(t, st, exec)
}

func listenProfileRoutes(t *testing.T, st store.Store, exec *session.Executor) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	h := server.New(server.WithListener(ln))
	h.GET("/v1/profile", httpserver.NewGetProfileHandler(st))
	h.PUT("/v1/profile", httpserver.NewPutProfileHandler(st))
	h.GET("/v1/sessions/:id", httpserver.NewGetSessionHandler(st))
	h.POST("/v1/sessions/end", httpserver.NewEndSessionHandler(st))
	if exec != nil {
		h.POST("/v1/sessions/stream", httpserver.NewStreamHandler(exec))
	}
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(50 * time.Millisecond)
	return "http://" + ln.Addr().String()
}

func getProfile(t *testing.T, baseURL, userID string) store.RelationshipProfile {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/profile?user_id="+userID, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d %s", resp.StatusCode, b)
	}
	var p store.RelationshipProfile
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatal(err)
	}
	return p
}

func putProfile(t *testing.T, baseURL, userID, jsonBody string) store.RelationshipProfile {
	t.Helper()
	req, err := http.NewRequest(http.MethodPut, baseURL+"/v1/profile?user_id="+userID, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d %s", resp.StatusCode, b)
	}
	var p store.RelationshipProfile
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatal(err)
	}
	return p
}
