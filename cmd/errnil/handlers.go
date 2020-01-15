package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"cpl.li/go/errnil/pkg/store"

	"github.com/gin-gonic/gin"
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
		if err != nil {
			if !errors.Is(err, store.ErrRepoNotFound) {
				log.Printf("failed getting repo(%s) from store, %s\n", repo, err.Error())
			}
		} else {
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
		}

		c.JSON(http.StatusOK, entry)
	}

}
