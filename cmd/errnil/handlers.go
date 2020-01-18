package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cpl.li/go/errnil/pkg/store"
)

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}

func handleInspect(downloadDir string, storage store.Store, cacheDuration time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get repo name
		repo := c.Query("repo")
		if repo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing repo query argument",
			})
			return
		}

		// check for cached entry
		entry, err := storage.GetEntry(repo)
		if err == nil {
			if time.Now().UTC().Sub(entry.UpdatedAt) < cacheDuration {
				c.JSON(http.StatusOK, entry)
				return
			}
		}

		// get and update entry
		entry, err = downloadInspectStore(repo, downloadDir, storage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}
