package disk

import (
	"crawlquery/pkg/domain"
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

type DiskRepository struct {
	filepath string
	jobs     []*domain.CrawlJob
	lock     sync.Mutex
}

func NewDiskRepository(filepath string) *DiskRepository {
	return &DiskRepository{
		filepath: filepath,
	}
}

func (dr *DiskRepository) Save() error {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	// Create a file for writing.
	file, err := os.Create(dr.filepath)
	if err != nil {
		return err
	}

	// Create a new gob encoder writing to the file.
	encoder := gob.NewEncoder(file)

	// Encode (serialize) the queue.
	if err := encoder.Encode(dr.jobs); err != nil {
		return err
	}

	return nil
}

func (dr *DiskRepository) Load() error {
	// Open the file for reading.
	file, err := os.Open(dr.filepath)

	if err != nil {
		return err
	}

	// Create a gob decoder
	decoder := gob.NewDecoder(file)

	// Create an empty Queue where the data will be decoded
	var jobs []*domain.CrawlJob

	if err := decoder.Decode(&jobs); err != nil {
		return err
	}

	dr.jobs = jobs
	return nil
}

func (dr *DiskRepository) Push(j *domain.CrawlJob) error {
	dr.jobs = append(dr.jobs, j)

	dr.Save()
	return nil
}

func (dr *DiskRepository) Pop() (*domain.CrawlJob, error) {
	if len(dr.jobs) == 0 {
		return nil, errors.New("no jobs in queue")
	}

	j := dr.jobs[0]
	dr.jobs = dr.jobs[1:]

	dr.Save()
	return j, nil
}
