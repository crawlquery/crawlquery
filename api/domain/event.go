package domain

type EventKey string
type Event interface {
	Key() EventKey
}

type EventService interface {
	Publish(event Event) error
	Subscribe(key EventKey, subFunc SubFunc)
}

type SubFunc func(event Event)
