package badger

import (
	"fmt"
	"net/url"
)

func ImageURL(repo, style string) string {
	format := "https://img.shields.io/endpoint?%s"

	if style == "" {
		style = badgeStyle
	}

	badgeURL := url.Values{}
	badgeURL.Set("repo", repo)
	badgeURL.Set("style", style)

	endpointURL := url.Values{}
	endpointURL.Set("url", fmt.Sprintf("https://errnil.cpl.li/api/badge?%s", badgeURL.Encode()))

	return fmt.Sprintf(format, endpointURL.Encode())
}

func ToMarkdown(repo, style string) string {
	format := "[![err != nil](%s)](https://errnil.cpl.li/api/inspect?repo=%s)"

	return fmt.Sprintf(format, ImageURL(repo, style), repo)
}
