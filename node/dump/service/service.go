package service

import "crawlquery/node/domain"

type Service struct {
	pageService    domain.PageService
	keywordService domain.KeywordService
}
