package http

import (
	"embed"
	"fmt"

	// "io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/utils/template"
)

func NewRouter(staticFiles *embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())
	// staticContent, _ := fs.Sub(staticFiles, "web/static")
	// router.StaticFS("/static", http.FS(staticContent))
	router.Static("/static", "web/static")
	router.Static("/novnc", "web/noVNC")
	router.StaticFS("/tools", http.Dir("web/component"))
	router.HTMLRender = template.LoadTemplates("web")
	return router
}

func NewServer(router *gin.Engine, configuration *environment.Configuration) error {
	// certFile := "cert/cert.pem"
	// keyFile := "cert/key.pem"
	// return router.RunTLS(fmt.Sprintf(":%s", configuration.Server.Port), certFile, keyFile)
	return router.Run(fmt.Sprintf(":%s", configuration.Server.Port))
}
