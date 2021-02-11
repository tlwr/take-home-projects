package link_parser

import (
	"net/url"
	"strings"

	"github.com/goware/urlx" // net/url url.Parse is not great for strict URL parsing
)

// Parse parses a link which may or may not be relative to the current URL
// Parse can return a nil value for URL at the same time as a nil error
func Parse(current *url.URL, link string) (u *url.URL, err error) {
	linkWithoutAnchor := strings.TrimSpace(strings.SplitN(link, "#", 2)[0])

	if strings.HasPrefix(linkWithoutAnchor, "mailto:") || len(linkWithoutAnchor) == 0 {
		// we cannot scrape mailto links for obvious reasons,
		//but they do not represent an error, just not a tangible link
		return nil, nil
	} else if strings.HasPrefix(linkWithoutAnchor, "/") {
		return parseHostRelative(current, linkWithoutAnchor), nil
	} else if strings.Contains(linkWithoutAnchor, ":") {
		return urlx.Parse(link)
	} else {
		return parseRelative(current, link), nil
	}
}

func parseHostRelative(current *url.URL, newPath string) *url.URL {
	// this copies the URL struct
	// this is kind of a hack because URL has a pointer to a UserInfo struct which is not copied
	// this is fine because we only care about modifying the Path
	var copyCurrent url.URL = *current
	copyCurrent.Path = newPath
	return &copyCurrent
}

func parseRelative(current *url.URL, newPath string) *url.URL {
	// this copies the URL struct
	// this is kind of a hack because URL has a pointer to a UserInfo struct which is not copied
	// this is fine because we only care about modifying the Path
	var copyCurrent url.URL = *current
	copyCurrent.Path += "/" + newPath
	return &copyCurrent
}
