package main

import (
	"fmt"
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

func handleBadge(storage store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get repo name
		repo := c.Query("repo")
		if repo == "" {
			c.Redirect(http.StatusTemporaryRedirect,
				fmt.Sprintf("%s?label=%s&message=missing repo&color=critical&style=%s",
					shieldsEndpoint,
					badgeLabel,
					c.Query("style"),
				))
			return
		}

		entry, err := storage.GetEntry(repo)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect,
				fmt.Sprintf("%s?label=%s&message=nil&color=inactive&style=%s",
					shieldsEndpoint,
					badgeLabel,
					c.Query("style"),
				))
			return
		}

		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s?label=%s&message=%d&color=%s&style=%s",
				shieldsEndpoint,
				badgeLabel,
				entry.PositionsCount,
				badgeColor,
				c.Query("style"),
			))
	}
}
