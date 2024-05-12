package domain

import "time"

type CrawlJob struct {
	ID        string
	URL       string
	CreatedAt time.Time
}
