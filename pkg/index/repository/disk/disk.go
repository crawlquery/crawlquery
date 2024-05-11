package disk

import (
	"crawlquery/pkg/index"
	"encoding/gob"
	"os"
)

type DiskRepository struct {
	filepath string
}

func NewDiskRepository(filepath string) *DiskRepository {
	return &DiskRepository{
		filepath: filepath,
	}
}

func (dr *DiskRepository) Save(idx *index.Index) error {
	// Create a file for writing.
	file, err := os.Create(dr.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob encoder writing to the file.
	encoder := gob.NewEncoder(file)

	// Encode (serialize) the index.
	if err := encoder.Encode(idx); err != nil {
		return err
	}

	return nil
}

func (dr *DiskRepository) Load() (*index.Index, error) {
	// Open the file for reading.
	file, err := os.Open(dr.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a gob decoder
	decoder := gob.NewDecoder(file)

	// Create an empty Index where the data will be decoded
	var idx index.Index
	if err := decoder.Decode(&idx); err != nil {
		return nil, err
	}
	return &idx, nil
}
