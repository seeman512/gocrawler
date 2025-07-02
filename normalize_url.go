package main

import (
	"net/url"
	"strings"
)

func normalizeURL(urlStr string) (string, error) {
	urlObj, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	// fmt.Printf("PATH %s: h: %s, p: %s, q: %s\n", urlStr, urlObj.Host, urlObj.Path, urlObj.RawQuery)
	if urlObj.RawQuery == "" {
		path := urlObj.Path
		if strings.HasSuffix(path, "/") {
			return urlObj.Host + path[:len(path)-1], nil
		}

		return urlObj.Host + path, nil
	}

	return urlObj.Host + urlObj.Path + "?" + urlObj.RawQuery, nil
}
