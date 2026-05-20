package api

import (
	"context"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterWebStatic serves the harness web/ assets. Hertz Static maps URL paths onto
// FS.Root verbatim; without PathRewrite, /js/app.js would resolve to Root/js/app.js
// (404 when Root is already web/js).
func RegisterWebStatic(r route.IRoutes, webRoot string) {
	r.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.File(filepath.Join(webRoot, "index.html"))
	})
	staticDir := func(urlPrefix, subdir string) {
		r.StaticFS(urlPrefix, &app.FS{
			Root:        filepath.Join(webRoot, subdir),
			Compress:    false,
			PathRewrite: app.NewPathSlashesStripper(1),
		})
	}
	staticDir("/css", "css")
	staticDir("/js", "js")
	r.StaticFile("/manifest.webmanifest", filepath.Join(webRoot, "manifest.webmanifest"))
}
