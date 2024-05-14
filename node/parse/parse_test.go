package parse_test

import (
	"crawlquery/node/parse"
	"fmt"
	"testing"

	testdataloader "github.com/peteole/testdata-loader"
)

func TestParse(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	res, err := parse.Parse(string(testdata), "http://google.com")

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	fmt.Printf("Parsed: %+v", res.Content)

	t.Fail()
}
