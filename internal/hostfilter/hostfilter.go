package hostfilter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gobwas/glob"
)

type HostFilter interface {
	IsAllowed(u *url.URL) bool
}

type hostFilter struct {
	globs []glob.Glob
}

func New(hosts []string) (HostFilter, error) {
	globs := []glob.Glob{}

	for _, h := range hosts {
		g, err := glob.Compile(strings.ToLower(h))
		if err != nil {
			return nil, fmt.Errorf("host glob (%s) is not a valid glob %w", h, err)
		}

		globs = append(globs, g)
	}

	return &hostFilter{
		globs: globs,
	}, nil
}

func (hf *hostFilter) IsAllowed(u *url.URL) bool {
	if u == nil || len(u.Host) == 0 {
		return false
	}

	for _, g := range hf.globs {
		if g.Match(strings.ToLower(u.Host)) {
			return true
		}
	}

	return false
}
