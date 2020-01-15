package main // import "cpl.li/go/errnil/cmd/errnil"

import (
	"log"
	"os"
	"path"
	"time"

	"cpl.li/go/errnil/pkg/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// setup download path
	downloadDir := os.Getenv("ERRNIL_DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = path.Join(os.TempDir(), "errnil")
	}

	storage := store.NewPrimitiveStore()

	// setup web server
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/health", handleHealth)
	router.Use(gin.Logger())

	router.GET("/inspect", handleInspect(downloadDir, storage, time.Minute))

	// start web server
	log.Fatal(router.Run(":8080"))
}
