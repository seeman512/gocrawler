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

	body, err := getHTML(url)
	if err != nil {
		fmt.Printf("Could not open page: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("BODY: %s\n", body)
}
