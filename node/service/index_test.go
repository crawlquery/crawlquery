package service_test

import (
	"bytes"
	"crawlquery/node/service"
	"crawlquery/pkg/index"
	"encoding/gob"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestSaveIndex(t *testing.T) {
	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "index-*.gob")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Close the temp file as SaveIndex will open it again for writing
	tmpFile.Close()

	// Create an instance of IndexService and a dummy Index
	service := service.NewIndexService()
	idx := &index.Index{
		// Initialize your Index with some data
		Forward: map[string]index.Document{
			"doc1": {ID: "doc1", Title: "Example Document"},
		},
		Inverted: map[string][]index.Posting{
			"example": {{DocumentID: "doc1", Frequency: 1, Positions: []int{0}}},
		},
	}

	// Call SaveIndex
	if err := service.SaveIndex(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	// Verify the file content
	content, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read back the saved file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("File is empty, but expected serialized index data")
	}

	// Optionally, you could deserialize the content back into an Index and compare with the original
	// This step is more complex and requires setting up a reader and using gob.Decoder
	var readIndex index.Index
	if err := deserializeIndex(content, &readIndex); err != nil {
		t.Errorf("Failed to deserialize index: %v", err)
	}

	if !reflect.DeepEqual(idx, &readIndex) {
		t.Errorf("Deserialized index does not match the original")
	}
}

// Helper function to deserialize index from bytes
func deserializeIndex(data []byte, idx *index.Index) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(idx)
}
