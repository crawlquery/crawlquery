package service

import (
	"crawlquery/api/domain"
)

type Service struct {
	subscribers map[domain.EventKey][]domain.SubFunc
}

func NewService() *Service {
	return &Service{
		subscribers: make(map[domain.EventKey][]domain.SubFunc),
	}
}

func (s *Service) Publish(event domain.Event) error {
	key := event.Key()
	for _, subFunc := range s.subscribers[key] {
		subFunc(event)
	}
	return nil
}

func (s *Service) Subscribe(key domain.EventKey, subFunc domain.SubFunc) {
	s.subscribers[key] = append(s.subscribers[key], subFunc)
}
