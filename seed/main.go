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
			panic(err)
		}

		if res.StatusCode != http.StatusCreated {
			panic(fmt.Errorf("unexpected status code: %d", res.StatusCode))
		}

		fmt.Printf("Crawling %s\n", domain)
	}

}
