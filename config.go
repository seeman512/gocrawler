package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

const maxPages = 50

type config struct {
	pages              map[string]int
	rawBaseURL         string
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func NewConfig(rawBaseURL string, maxConcurrent int) (*config, error) {

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	return &config{
		pages:              map[string]int{},
		rawBaseURL:         rawBaseURL,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrent),
		wg:                 &sync.WaitGroup{},
	}, nil
}

func (cfg *config) crawlPage(rawCurrentURL string) error {

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

	fmt.Printf("Start url: %s\n", rawCurrentURL)

	if len(cfg.pages) > maxPages {
		return fmt.Errorf("pages limit exceed")
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

		if cfg.addPageVisit(nUrl) {
			cfg.wg.Add(1)
			fmt.Printf("Start: %s\n", nUrl)
			cfg.concurrencyControl <- struct{}{}
			go func() {
				defer func() {
					cfg.wg.Done()
					fmt.Printf("Stop: %s\n", nUrl)
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
