package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	var seedFile string

	flag.StringVar(&seedFile, "file", "global/domains.txt", "File containing domains to seed")

	flag.Parse()

	file, err := os.ReadFile(seedFile)

	if err != nil {
		panic(err)
	}

	splitByLine := strings.Split(string(file), "\n")

	for _, domain := range splitByLine {
		res, err := http.Post(
			"http://localhost:8080/pages",
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, domain))))

		if err != nil {
			continue
		}

		if res.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to crawl %s got unexpected status code: %d", domain, res.StatusCode)
		}

		fmt.Printf("Crawling %s\n", domain)
	}

}
