package event_bus

import "context"

type IEventBus interface {
	Register(sub IEventSubscriber, events ...IEvent)
	UnRegister(sub IEventSubscriber, events ...IEvent)
	Publish(events ...IEventIns)
}

type IEvent interface {
	EventName() string
}

type IEventIns interface {
	IEvent
	Context() context.Context
}

type IEventHandle interface {
	IEvent
	Handle(event IEventIns)
}

type IEventPublisher interface {
	Init(bus IEventBus)
	PubEvent(events ...IEventIns)
}

type IEventSubscriber interface {
	Init(name string, bus IEventBus, events ...IEventHandle)
	Subscriber() string
	Subscribe(events ...IEventHandle)
	UnSubscribe(events ...IEvent)
	OnEvent(event IEventIns)
}
