package main // import "cpl.li/go/errnil/cmd/errnil"

import (
	"log"
	"os"
	"path"
	"time"

	"cpl.li/go/errnil/pkg/store"

	"github.com/gin-gonic/gin"
)

const (
	shieldsEndpoint = "https://img.shields.io/static/v1"
	badgeColor      = "e44"
	badgeLabel      = "err != nil"
)

func main() {
	downloadDir := path.Join(os.TempDir(), "errnil")
	storage := store.NewPrimitiveStore()

	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.GET("/health", handleHealth)
		api.GET("/inspect", handleInspect(downloadDir, storage, time.Minute))
		api.GET("/badge", handleBadge(storage))
	}

	log.Fatal(router.Run(":8080"))
}
