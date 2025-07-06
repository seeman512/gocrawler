// usage: ./crawler URL maxConcurrency maxPages

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func exitError(err error) {
	fmt.Printf("Could not create config: %v\n", err)
	os.Exit(1)
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
	cfg.showReport()
	// fmt.Printf("PAGES FINAL: %v\n", cfg.pages)
}
