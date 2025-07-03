package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) error {
	baseObj, err := url.Parse(rawBaseURL)
	if err != nil {
		return err
	}

	currentObj, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}

	if baseObj.Host != currentObj.Host {
		return fmt.Errorf("base host: %s differ current host: %s", baseObj.Host, currentObj.Host)
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}

	pages[normCurrentURL] = 1

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return err
	}

	urls, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		return err
	}

	if len(pages) > 50 {
		return nil
	}

	for _, u := range urls {

		urlObj, err := url.Parse(u)
		if err != nil {
			fmt.Printf("Invalid url: %s\n", u)
			continue
		}

		if baseObj.Host != urlObj.Host {
			fmt.Printf("base host: %s differ current host: %s", baseObj.Host, urlObj.Host)
			continue
		}

		nUrl, err := normalizeURL(u)
		if err != nil {
			fmt.Printf("Could not normilize url: %s\n", u)
			continue
		}

		n, ok := pages[nUrl]
		if ok {
			pages[nUrl] = n + 1
			continue
		}

		fmt.Printf("PAGES: %v\n", pages)
		time.Sleep(2 * time.Second)
		pages[nUrl] = 1
		crawlPage(rawBaseURL, u, pages)
	}

	return nil
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode > 400 {
		return "", fmt.Errorf("wrong page: %d", resp.StatusCode)
	}

	contentType := strings.ToLower(resp.Header.Get("content-type"))
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("not html page: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)

	if err != nil {
		return nil, err
	}

	links := []string{}

	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, a := range node.Attr {
				if a.Key == "href" {

					if strings.HasPrefix(a.Val, "http") {
						links = append(links, a.Val)
					} else {
						links = append(links, rawBaseURL+a.Val)
					}
					break
				}
			}
		}
	}

	fmt.Printf("LINKS: %v\n", links)

	return links, nil
}
