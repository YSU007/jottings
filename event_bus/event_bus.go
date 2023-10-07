package event_bus

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]map[string]EventSubscribe),
	}
}

type EventBus struct {
	subscribers map[string]map[string]EventSubscribe
}

func (eb *EventBus) Register(sub EventSubscribe, events ...Event) {
	for _, event := range events {
		subs := eb.subscribers[event.Name()]
		if subs == nil {
			subs = make(map[string]EventSubscribe)
		}
		subs[sub.Subscriber()] = sub
		eb.subscribers[event.Name()] = subs
	}
}

func (eb *EventBus) UnRegister(sub EventSubscribe, events ...Event) {
	for _, event := range events {
		subs := eb.subscribers[event.Name()]
		if subs != nil {
			delete(subs, sub.Subscriber())
		}
	}
}

func (eb *EventBus) Publish(event Event) {
	for _, sub := range eb.subscribers[event.Name()] {
		sub.OnEvent(event)
	}
}
