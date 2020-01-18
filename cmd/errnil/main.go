package main // import "cpl.li/go/errnil/cmd/errnil"

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"

	"cpl.li/go/errnil/pkg/store"
)

func main() {
	downloadDir := env("ERRNIL_DOWNLOAD_DIR", path.Join(os.TempDir(), "errnil"))
	storage := store.NewPrimitiveStore()

	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.GET("/health", handleHealth)
		api.GET("/inspect", handleInspect(downloadDir, storage, time.Minute))
		api.GET("/badge", handleBadge(storage))
	}

	log.Fatal(router.Run(env("ERRNIL_ADDRESS", ":8080")))
}
