package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/internal/configpaths"
	"github.com/huodaoshi/harness/backend/internal/httpserver"
	"github.com/huodaoshi/harness/backend/internal/session"
)

func main() {
	ctx := context.Background()
	exec, err := session.NewExecutor(ctx)
	if err != nil {
		log.Fatalf("session executor: %v", err)
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	h := server.Default(server.WithHostPorts(addr))
	h.POST("/v1/sessions/stream", httpserver.NewStreamHandler(exec))
	h.GET("/v1/profile", httpserver.NewGetProfileHandler(exec.Store))
	h.PUT("/v1/profile", httpserver.NewPutProfileHandler(exec.Store))

	webRoot := configpaths.WebRoot()
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.File(filepath.Join(webRoot, "index.html"))
	})
	h.Static("/css", filepath.Join(webRoot, "css"))
	h.Static("/js", filepath.Join(webRoot, "js"))
	h.StaticFile("/manifest.webmanifest", filepath.Join(webRoot, "manifest.webmanifest"))

	log.Printf("listening on %s (web=%s)", addr, webRoot)
	h.Spin()
}
