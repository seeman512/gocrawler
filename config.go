package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type config struct {
	maxPages           int
	pages              map[string]int
	rawBaseURL         string
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func NewConfig(rawBaseURL string, maxConcurrent, maxPages int) (*config, error) {

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Start config: %d => %d\n", maxConcurrent, maxPages)

	return &config{
		maxPages:           maxPages,
		pages:              map[string]int{},
		rawBaseURL:         rawBaseURL,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrent),
		wg:                 &sync.WaitGroup{},
	}, nil
}

func (cfg *config) crawlPage(rawCurrentURL string) error {
	fmt.Printf("Start url: %s\n", rawCurrentURL)

	if cfg.pagesLimitExceed() {
		return fmt.Errorf("pages limit exceed")
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}

	if cfg.baseURL.Host != currentURL.Host {
		return fmt.Errorf("base host: %s differ current host: %s", cfg.baseURL.Host, currentURL.Host)
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}

	cfg.mu.Lock()
	if _, ok := cfg.pages[normCurrentURL]; !ok {
		cfg.pages[normCurrentURL] = 1
	} else {
		cfg.mu.Unlock()
		return nil
	}

	cfg.mu.Unlock()

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return err
	}

	urls, err := getURLsFromHTML(html, cfg.rawBaseURL)
	if err != nil {
		return err
	}

	for _, u := range urls {

		currURL, err := url.Parse(u)
		if err != nil {
			fmt.Printf("Invalid url: %s\n", u)
			continue
		}

		if cfg.baseURL.Host != currURL.Host {
			fmt.Printf("base host: %s differ current host: %s\n", cfg.baseURL.Host, currURL.Host)
			continue
		}

		nUrl, err := normalizeURL(u)
		if err != nil {
			fmt.Printf("Could not normilize url: %s\n", u)
			continue
		}

		if cfg.pagesLimitExceed() {
			return fmt.Errorf("pages limit exceed")
		}

		if cfg.addPageVisit(nUrl) {
			cfg.wg.Add(1)
			cfg.concurrencyControl <- struct{}{}
			go func() {
				defer func() {
					cfg.wg.Done()
				}()

				time.Sleep(1 * time.Second)
				cfg.crawlPage(u)

				<-cfg.concurrencyControl
			}()
		}
	}

	return nil
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	n, ok := cfg.pages[normalizedURL]
	cfg.pages[normalizedURL] = n + 1
	return !ok
}

func (cfg *config) pagesLimitExceed() bool {

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages) >= cfg.maxPages
}

func (cfg *config) showReport() {
	for page, cnt := range cfg.pages {
		fmt.Printf("Page: %s => %d\n", page, cnt)
	}
}
