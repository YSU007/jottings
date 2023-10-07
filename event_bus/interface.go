package event_bus

type Event interface {
	Name() string
	IntParam() []int64
	StringParam() []string
}

type EventPublish interface {
	PubEvent(event Event)
}

type EventSubscribe interface {
	Subscriber() string
	OnEvent(event Event)
}
