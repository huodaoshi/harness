package main

import (
	"context"
	"log"
	"os"

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
	h.GET("/v1/sessions/:id", httpserver.NewGetSessionHandler(exec.Store))
	h.POST("/v1/sessions/end", httpserver.NewEndSessionHandler(exec.Store))
	h.GET("/v1/profile", httpserver.NewGetProfileHandler(exec.Store))
	h.PUT("/v1/profile", httpserver.NewPutProfileHandler(exec.Store))

	webRoot := configpaths.WebRoot()
	httpserver.RegisterWebStatic(h, webRoot)

	log.Printf("listening on %s (web=%s)", addr, webRoot)
	h.Spin()
}
