package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/event/service"
	"testing"
)

type testEvent struct {
}

func (e testEvent) Key() domain.EventKey {
	return "test"
}

func TestEventService(t *testing.T) {
	t.Run("can publish and subscribe to events", func(t *testing.T) {
		service := service.NewService()

		event := testEvent{}

		var receivedEvent domain.Event

		subFunc := func(e domain.Event) {
			receivedEvent = e
		}

		key := domain.EventKey("test")

		service.Subscribe(key, subFunc)

		err := service.Publish(event)
		if err != nil {
			t.Errorf("Error publishing event: %v", err)
		}

		if receivedEvent != event {
			t.Errorf("Expected received event to be the same as published event")
		}
	})
}
