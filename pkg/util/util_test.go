package util_test

import (
	"crawlquery/pkg/util"
	"testing"
)

func TestValidatePageID(t *testing.T) {

	invalidPageIDs := []string{
		"page1$",
		"/page1",
		"page1/",
		"(page1)",
		"page1?do2omndomdodod",
		"page1#",
		"page1&",
		"page1=",
		"page1%",
		"page1+",
		"page1-",
		"p@age1_",
		// 32 characters
		"page1page1page1page1page!page1pa",
	}

	for _, invalidPageID := range invalidPageIDs {

		if util.ValidatePageID(invalidPageID) {
			t.Fatalf("Test failed: Invalid page ID was accepted.")
		}
	}

	validPageIDs := []string{
		"jdjdjdjdjdjdjdjdjdjdjdjdjdjdjdj",
		"3djdjdjdjd2djdjdjdjdjdjdjdjdjdj",
	}

	for _, validPageID := range validPageIDs {
		if !util.ValidatePageID(validPageID) {
			t.Fatalf("Test failed: Valid page ID was rejected.")
		}
	}
}

// TestPageID is a simple test function for PageID.
func TestPageID(t *testing.T) {
	url1 := "https://example.com"
	url2 := "https://example.com/page"
	url3 := "https://example.org"

	hash1 := util.PageID(url1)
	hash2 := util.PageID(url2)
	hash3 := util.PageID(url3)

	// Check for uniqueness
	if hash1 == hash2 || hash1 == hash3 || hash2 == hash3 {
		t.Fatalf("Test failed: Hashes are not unique.")
	}

	if len(hash1) != 32 || len(hash2) != 32 || len(hash3) != 32 {
		t.Fatalf("Test failed: Hashes are not 32 characters long.")
	}

	// only contains alphanumeric characters
	for _, c := range hash1 {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Fatalf("Test failed: Hash contains non-alphanumeric characters.")
		}
	}

	for _, c := range hash2 {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Fatalf("Test failed: Hash contains non-alphanumeric characters.")
		}
	}

	for _, c := range hash3 {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Fatalf("Test failed: Hash contains non-alphanumeric characters.")
		}
	}
}
