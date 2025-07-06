package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
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

	// fmt.Printf("LINKS: %v\n", links)

	return links, nil
}
