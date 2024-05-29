package service

import "crawlquery/api/domain"

type Service struct {
	crawlJobRepo domain.CrawlJobRepository
	crawlLogRepo domain.CrawlLogRepository
}
