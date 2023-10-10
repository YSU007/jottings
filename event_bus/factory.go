package event_bus

type EventBusType int

const (
	BusSingle EventBusType = iota
	BusSafety
)

func NewEventBus(t EventBusType) IEventBus {
	switch t {
	case BusSingle:
		return NewEventBusSingle()
	case BusSafety:
		const bucketNum = 32
		return NewEventBusBucket(bucketNum)
	}
	return nil
}

type EventSubscriberType int

const (
	SubscriberSingle EventSubscriberType = iota
	SubscriberSafety
)

func NewEventSubscriber(t EventSubscriberType, name string, bus IEventBus, events ...IEventHandle) IEventSubscriber {
	switch t {
	case SubscriberSingle:
		return NewEventSubscriberSingle(name, bus, events...)
	case SubscriberSafety:
		return NewEventSubscriberSafety(name, bus, events...)
	}
	return nil
}
