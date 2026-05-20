package apitest

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/redis/go-redis/v9"

	"github.com/huodaoshi/harness/backend/api"
	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/chatmodel"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

func TestStreamHandler_RateLimit429(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	limiter := infra.NewRedisRateLimiter(rdb)

	ctx := context.Background()
	exec, err := application.NewExecutorWithGateway(ctx, store.NewMemoryStore(), chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}

	const limit = 2
	baseURL := newStreamServerWithLimiter(t, exec, limiter, limit)
	body := `{"message":"hi","mode":"normal"}`

	for i := 0; i < limit; i++ {
		resp := postStream(t, baseURL, body)
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("request %d: status=%d body=%s", i+1, resp.StatusCode, b)
		}
		resp.Body.Close()
	}

	resp := postStream(t, baseURL, body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTooManyRequests {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d want 429 body=%s", resp.StatusCode, b)
	}
	var er struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		t.Fatal(err)
	}
	if er.Code != 4013 || er.Message == "" {
		t.Fatalf("payload: %+v", er)
	}
}

func TestStreamHandler_RateLimitDoesNotBlockCrisis(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	limiter := infra.NewRedisRateLimiter(rdb)

	ctx := context.Background()
	exec, err := application.NewExecutorWithGateway(ctx, store.NewMemoryStore(), chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}

	baseURL := newStreamServerWithLimiter(t, exec, limiter, 1)
	body := `{"message":"我不想活了","mode":"distress"}`

	resp := postStream(t, baseURL, body)
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("first crisis: status=%d %s", resp.StatusCode, b)
	}
	assertSSEHasEvents(t, resp.Body, sseExpect{wantCrisis: true})
	resp.Body.Close()

	resp2 := postStream(t, baseURL, body)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("second crisis: status=%d want 429", resp2.StatusCode)
	}
	if exec.ChatCalls.Load() != 0 {
		t.Fatalf("chat calls=%d want 0", exec.ChatCalls.Load())
	}
}

func newStreamServerWithLimiter(t *testing.T, exec *application.Executor, limiter *infra.RedisRateLimiter, perMinute int) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	guestMw := testGuestMw(t)
	h := server.New(server.WithListener(ln))
	h.POST("/v1/sessions/stream", guestMw, api.NewStreamHandler(exec, limiter, perMinute))
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(50 * time.Millisecond)
	return "http://" + ln.Addr().String()
}
