package main // import "cpl.li/go/errnil/cmd/errnil"

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cpl.li/go/errnil/pkg/errnil"
)

func main() {
	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	dir := path.Join(os.TempDir(), "errnil")
	router.GET("/inspect", func(c *gin.Context) {
		pkg := c.Query("pkg")

		downloadPath, err := errnil.Download(pkg, dir)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer os.RemoveAll(downloadPath)

		positions, err := errnil.Inspect(downloadPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		cleanPositions(positions, path.Join(downloadPath, ".."))

		c.JSON(http.StatusOK, gin.H{
			"positions_count": len(positions),
			"positions":       positions,
		})
	})

	log.Fatal(router.Run(":8080"))
}

func cleanPositions(positions []errnil.Position, tmpDir string) {
	for idx := range positions {
		positions[idx].Filename = strings.Replace(positions[idx].Filename, tmpDir+"/", "", 1)
	}
}
