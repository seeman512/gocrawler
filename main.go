// usage: ./crawler URL maxConcurrency maxPages

package main

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func exitError(err error) {
	fmt.Printf("Could not create config: %v\n", err)
	os.Exit(1)
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`=============================
  REPORT for %s
=============================
`, baseURL)

	type item struct {
		page string
		cnt  int
	}

	pagesList := make([]item, len(pages))
	i := 0
	for k, v := range pages {
		pagesList[i] = item{k, v}
		i++
	}

	slices.SortFunc(pagesList, func(p1, p2 item) int {
		return cmp.Or(
			-cmp.Compare(p1.cnt, p2.cnt),
			cmp.Compare(p1.page, p2.page),
		)
	})

	for _, p := range pagesList {
		fmt.Printf("Found %d internal links to %s\n", p.cnt, p.page)
	}
}

func main() {
	if len(os.Args) < 4 {
		exitError(errors.New("no website,maxConcurrency, maxPages provided"))
	}

	if len(os.Args) > 4 {
		exitError(errors.New("too many arguments provided"))
	}

	url := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", url)

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		exitError(err)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		exitError(err)
	}

	cfg, err := NewConfig(url, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Could not create config: %v\n", err)
		os.Exit(1)
	}

	err = cfg.crawlPage(url)
	if err != nil {
		fmt.Printf("Stop crawl pages: %v\n", err)
	}

	cfg.wg.Wait()
	printReport(cfg.pages, url)
}
