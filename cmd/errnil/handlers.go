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
		"timestamp": time.Now().UTC().Unix(),
	})
}

func handleGetEntry(storage store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := c.Query("repo")
		if repo == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "missing repo query argument",
			})
			return
		}

		entry, err := storage.GetEntry(repo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "entry not found in storage, " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}

func handleInspect(downloadDir string, storage store.Store, cacheDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := c.Query("repo")
		if repo == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "missing repo query argument",
			})
			return
		}

		entry, err := storage.GetEntry(repo)
		if err == nil {
			if time.Now().UTC().Sub(entry.UpdatedAt) < cacheDuration {
				c.Next()
				return
			}
		}

		_, err = downloadInspectStore(repo, downloadDir, storage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()
	}
}
