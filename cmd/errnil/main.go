package main // import "cpl.li/go/errnil/cmd/errnil"

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	gintemplate "github.com/foolin/gin-template"

	"github.com/gin-gonic/gin"

	"cpl.li/go/errnil/pkg/badger"
	"cpl.li/go/errnil/pkg/store"
)

func main() {
	downloadDir := env("ERRNIL_DOWNLOAD_DIR", path.Join(os.TempDir(), "errnil"))
	storage := store.NewPrimitiveStore()

	router := gin.New()
	router.Use(gin.Recovery())

	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         path.Join("static", "views"),
		Extension:    ".gohtml",
		Partials:     []string{"layouts/reporesults"},
		Funcs:        template.FuncMap{},
		DisableCache: true,
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})
	router.GET("/results", func(c *gin.Context) {
		repo := c.Query("repo")
		if repo == "" {
			c.HTML(http.StatusBadRequest, "error", gin.H{
				"code":    http.StatusBadRequest,
				"message": "missing repo query param",
			})
			return
		}

		entry, err := downloadInspectStore(repo, downloadDir, storage)
		if err != nil {
			if errors.Is(err, store.ErrRepoNotFound) {
				c.HTML(http.StatusNotFound, "error", gin.H{
					"code":    http.StatusNotFound,
					"message": err.Error(),
				})
			} else {
				c.HTML(http.StatusInternalServerError, "error", gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error(),
				})
			}
			return
		}

		c.HTML(http.StatusOK, "results", gin.H{
			"repo":          repo,
			"positionCount": entry.PositionsCount,
			"positions":     entry.Positions,
			"updatedAt":     entry.UpdatedAt,
			"markdown":      badger.ToMarkdown(repo, c.Query("style")),
			"imageURL":      badger.ImageURL(repo, c.Query("style")),
		})
	})

	api := router.Group("/api")
	{
		api.GET("/health", handleHealth)
		api.GET("/inspect", handleInspect(downloadDir, storage, time.Minute*5), handleGetEntry(storage))
		api.GET("/badge", handleInspect(downloadDir, storage, time.Minute*5), badger.Badge(storage))
	}

	log.Fatal(router.Run(
		env("ERRNIL_ADDRESS", "") + ":" + env("PORT", "8080")))
}
