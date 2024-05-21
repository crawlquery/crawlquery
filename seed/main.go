package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	file, err := os.ReadFile("domains.txt")

	if err != nil {
		panic(err)
	}

	splitByLine := strings.Split(string(file), "\n")

	for _, domain := range splitByLine {
		res, err := http.Post(
			"http://localhost:8080/crawl-jobs",
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "https://%s"}`, domain))))

		if err != nil {
			continue
		}

		if res.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to crawl %s got unexpected status code: %d", domain, res.StatusCode)
		}

		fmt.Printf("Crawling %s\n", domain)
	}

}
