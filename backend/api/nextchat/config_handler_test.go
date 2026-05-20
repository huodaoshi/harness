package nextchat

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

	"github.com/huodaoshi/harness/backend/conf"
)

func TestConfigHandler_NoCode(t *testing.T) {
	backendRoot := filepath.Join("..", "..")
	if err := os.Chdir(backendRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(filepath.Join(backendRoot, "api", "nextchat")) })

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
	Register(h, c)
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
	var body dangerConfig
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if body.NeedCode {
		t.Fatal("expected needCode=false")
	}
	if !body.HideUserApiKey {
		t.Fatal("expected hideUserApiKey=true")
	}
	if body.CustomModels == "" {
		t.Fatal("expected customModels")
	}
}
