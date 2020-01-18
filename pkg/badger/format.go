package badger

import (
	"fmt"
	"net/url"
)

func ToMarkdown(repo, style string) string {
	format := "[![err != nil](https://img.shields.io/endpoint?%s)](https://errnil.cpl.li/api/inspect?repo=%s)"

	badgeURL := url.Values{}
	badgeURL.Set("repo", repo)
	badgeURL.Set("style", style)

	endpointURL := url.Values{}
	endpointURL.Set("url", fmt.Sprintf("https://errnil.cpl.li/api/badge?%s", badgeURL.Encode()))

	return fmt.Sprintf(format, endpointURL.Encode(), repo)
}
