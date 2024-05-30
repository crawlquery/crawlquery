package service

import (
	"crawlquery/api/domain"
	"net/url"
	"sync"
	"time"
)

type Service struct {
	domains map[string]time.Time
	mutex   *sync.Mutex
}

type Option func(*Service)

func WithRateLimit(rateLimit time.Duration) Option {
	return func(t *Service) {
		go func() {
			for {
				t.mutex.Lock()
				for k, v := range t.domains {
					if time.Since(v) > rateLimit {
						delete(t.domains, k)
					}
				}
				t.mutex.Unlock()
			}
		}()
	}
}

func NewService(opts ...Option) *Service {
	s := &Service{
		map[string]time.Time{},
		&sync.Mutex{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (t *Service) CheckAndThrottle(rawURL domain.URL) (bool, error) {
	parsed, err := url.ParseRequestURI(string(rawURL))

	if err != nil {
		return false, err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, ok := t.domains[parsed.Host]; ok {
		return false, nil
	}

	t.domains[parsed.Host] = time.Now()
	return true, nil
}
