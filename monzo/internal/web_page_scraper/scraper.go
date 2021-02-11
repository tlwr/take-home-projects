package web_page_scraper

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"

	"github.com/antchfx/htmlquery"

	link "github.com/tlwr/take-home-projects/monzo/internal/link_parser"
)

type ScrapeResult struct {
	URL         *url.URL
	Links       []*url.URL
	ParseErrors []error
}

type Scraper interface {
	Scrape(u *url.URL) (*ScrapeResult, error)
}

type webScraper struct {
	client *http.Client
}

func (c *webScraper) Scrape(u *url.URL) (*ScrapeResult, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "monzo-take-home-crawler")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("expected 200-399 status code received %d (%s)", resp.StatusCode, u.String())
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(node, "//a")
	if err != nil {
		return nil, err
	}

	links := []*url.URL{}
	errs := []error{}
	for _, node := range nodes {
		href := ""

		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}

		if href == "" {
			continue
		}

		l, err := link.Parse(u, href)
		if err != nil {
			errs = append(errs, err)
		} else if l != nil {
			links = append(links, l)
		}
	}

	return &ScrapeResult{
		URL:         u,
		Links:       links,
		ParseErrors: errs,
	}, nil
}

func New(h *http.Client) Scraper {
	return &webScraper{
		client: h,
	}
}
