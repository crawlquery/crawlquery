package index_test

import (
	"crawlquery/pkg/index"
	"testing"
)

func TestAddDocument(t *testing.T) {
	// Initialize index and document as before
	idx := index.NewIndex()
	doc := index.Document{
		ID:              "doc1",
		URL:             "http://example.com",
		Title:           "Test Document",
		Content:         `<html><head><title>Example</title></head><body><h1>Hello World!</h1><p>This is a simple test. Numbers: 1234.</p></body></html>`,
		MetaDescription: "A simple test document",
	}

	// Add document to the index
	idx.AddDocument(doc)

	// Retrieve the document from the forward index
	indexedDoc, exists := idx.Forward[doc.ID]
	if !exists {
		t.Fatalf("Document with ID %s not found in forward index", doc.ID)
	}

	// Compare each field
	if indexedDoc.ID != doc.ID ||
		indexedDoc.URL != doc.URL ||
		indexedDoc.Title != doc.Title ||
		indexedDoc.Content != doc.Content ||
		indexedDoc.MetaDescription != doc.MetaDescription {
		t.Errorf("Document fields do not match. Got %+v, want %+v", indexedDoc, doc)
	}
}
