package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/goware/urlx" // net/url url.Parse is not great for strict URL parsing

	fllag "github.com/tlwr/take-home-projects/monzo/internal/flag"
	"github.com/tlwr/take-home-projects/monzo/internal/hostfilter"
	queue "github.com/tlwr/take-home-projects/monzo/internal/url_dedup_queue"
	scraper "github.com/tlwr/take-home-projects/monzo/internal/web_page_scraper"

	cl "github.com/tlwr/operator-tools/pkg/colour"
)

type result struct {
	url   *url.URL
	links []*url.URL
}

func main() {
	var (
		u        string
		hosts    fllag.StringSliceFlag
		parallel int
	)

	flag.IntVar(&parallel, "parallel", 8, "parallelism factor")
	flag.StringVar(&u, "url", "", "page on which to begin crawling")
	flag.Var(&hosts, "host", "crawls pages from this host (valid multiple times) (can be a glob)")
	flag.Parse()

	if parallel < 1 || parallel > 256 {
		log.Fatal("parallel flag should be between 1 and 256")
	}

	if len(hosts) == 0 {
		log.Fatal("host flag is required")
	}

	if u == "" {
		log.Fatal("url flag is required")
	}

	ur, err := urlx.Parse(u)
	if err != nil {
		log.Fatalf("could not validate url: %v", err)
	}

	hf, err := hostfilter.New(hosts)
	if err != nil {
		log.Fatalf("could not construct a hostfilter from host flags: %v", err)
	}

	for _, host := range hosts {
		log.Printf("will include URLs within %s", host)
	}

	errC := make(chan error, 8)
	resultsC := make(chan result, 64)

	go func() {
		for err := range errC {
			log.Println(cl.Red(fmt.Sprintf("scraper encountered error: %s", err)))
		}
	}()

	resultsWg := sync.WaitGroup{}
	resultsWg.Add(1)
	go func() {
		for r := range resultsC {
			log.Printf("results for page: %s", cl.Green(r.url.String()))
			for _, l := range r.links {
				log.Println(cl.Blue(l.String()))
			}
		}
		resultsWg.Done()
	}()

	q := queue.New()
	s := scraper.New(http.DefaultClient)

	log.Printf("queueing first url %s", ur)
	q.Enqueue(ur)

	for wi := 0; wi <= parallel; wi++ {
		go q.Iter(func(u *url.URL) {
			res, err := s.Scrape(u)

			if err != nil {
				errC <- err
				return
			}

			for _, err := range res.ParseErrors {
				errC <- err
			}

			for _, link := range res.Links {
				if hf.IsAllowed(link) {
					q.Enqueue(link)
				}
			}

			resultsC <- result{
				url:   u,
				links: res.Links,
			}
		})
	}

	q.Wait()

	close(errC)
	close(resultsC)

	resultsWg.Wait()
}
