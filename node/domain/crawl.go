package domain

import "errors"

var ErrCrawlFailedToStoreHtml = errors.New("failed to store html")
var ErrCrawlFailedToFetchHtml = errors.New("failed to fetch html")
