package badger

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cpl.li/go/errnil/pkg/store"

	"github.com/gin-gonic/gin"
)

const (
	badgeColor = "e44"
	badgeLabel = "err != nil"
	badgeStyle = "flat"
)

type badge struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
	Style         string `json:"style"`
	IsError       bool   `json:"isError"`
}

func (b badge) JSON() string {
	data, _ := json.Marshal(&b)
	return string(data)
}

func newBadge(message, color, style string, isErr bool) badge {
	if strings.TrimSpace(style) == "" {
		style = badgeStyle
	}

	return badge{
		SchemaVersion: 1,
		Label:         badgeLabel,
		Message:       message,
		Color:         color,
		Style:         style,
		IsError:       isErr,
	}
}

// Badge builds a shields.io endpoint for generating errnil counter badges.
func Badge(storage store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := c.Query("repo")
		style := c.Query("style")

		if repo == "" {
			c.JSON(http.StatusOK,
				newBadge("missing repo", "critical", style, true))
			return
		}

		entry, err := storage.GetEntry(repo)
		if err != nil {
			c.JSON(http.StatusOK,
				newBadge("?", "inactive", style, true))
			return
		}

		c.JSON(http.StatusOK,
			newBadge(strconv.Itoa(entry.PositionsCount), badgeColor, c.Query("style"), false))
	}
}
