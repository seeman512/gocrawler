package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		println("too many arguments provided")
		os.Exit(1)
	}

	url := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", url)

	// body, err := getHTML(url)
	// pages := map[string]int{}
	// err := crawlPage(url, url, pages)

	cfg, err := NewConfig(url, 2)
	if err != nil {
		fmt.Printf("Could not create config: %v\n", err)
		os.Exit(1)
	}

	err = cfg.crawlPage(url)
	if err != nil {
		fmt.Printf("Could not crawl page: %v\n", err)
		os.Exit(1)
	}

	cfg.wg.Wait()

	fmt.Printf("PAGES FINAL: %v\n", cfg.pages)
}
