package nextchat_test

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/api/nextchat"
	"github.com/huodaoshi/harness/backend/conf"
)

func TestConfigHandler_NoCode(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	backendRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if err := os.Chdir(backendRoot); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CODE", "")
	t.Setenv("ARK_API_KEY", "test-key")
	t.Setenv("ARK_MODEL_ID", "ep-test")

	c, err := conf.Load()
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	h := server.New(server.WithListener(ln))
	nextchat.Register(h, c)
	go h.Spin()
	t.Cleanup(func() { _ = h.Close() })
	time.Sleep(30 * time.Millisecond)

	resp, err := http.Get("http://" + ln.Addr().String() + "/api/config")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	raw, _ := io.ReadAll(resp.Body)
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if body["needCode"] != false {
		t.Fatalf("expected needCode=false, got %v", body["needCode"])
	}
	if body["hideUserApiKey"] != true {
		t.Fatalf("expected hideUserApiKey=true, got %v", body["hideUserApiKey"])
	}
	if body["customModels"] == "" {
		t.Fatal("expected customModels")
	}
}
