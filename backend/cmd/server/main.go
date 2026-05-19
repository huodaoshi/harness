package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/internal/httpserver"
	"github.com/huodaoshi/harness/backend/internal/session"
)

func main() {
	ctx := context.Background()
	runnable, err := session.CompileDefaultGraph(ctx)
	if err != nil {
		log.Fatalf("compile graph: %v", err)
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	h := server.Default(server.WithHostPorts(addr))
	h.POST("/v1/sessions/stream", httpserver.NewStreamHandler(runnable))

	log.Printf("listening on %s", addr)
	h.Spin()
}
