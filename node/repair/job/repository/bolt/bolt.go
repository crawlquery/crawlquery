package bolt

import (
	"crawlquery/node/domain"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

type Repository struct {
	db *bolt.DB
}

var repairJobBucket = []byte("RepairJobs")

func NewRepository(db *bolt.DB) *Repository {

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(repairJobBucket)

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(job *domain.RepairJob) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(repairJobBucket)

		serialised, err := json.Marshal(job)

		if err != nil {
			return fmt.Errorf("serialise job: %s", err)
		}

		return b.Put([]byte(job.PageID), serialised)
	})
}

func (r *Repository) Get(pageID string) (*domain.RepairJob, error) {
	var job domain.RepairJob

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(repairJobBucket)

		v := b.Get([]byte(pageID))

		if v == nil {
			return fmt.Errorf("job not found")
		}

		return json.Unmarshal(v, &job)
	})

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *Repository) Update(job *domain.RepairJob) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(repairJobBucket)

		serialised, err := json.Marshal(job)

		if err != nil {
			return fmt.Errorf("serialise job: %s", err)
		}

		return b.Put([]byte(job.PageID), serialised)
	})
}
