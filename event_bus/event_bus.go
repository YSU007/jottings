package event_bus

import (
	"context"
	"hash/fnv"
	"sync"
)

type Event struct {
	name string
}

func (e *Event) EventName() string {
	return e.name
}

func NewEventIns(e Event, ctx context.Context) IEventIns {
	return &EventIns{
		Event: e,
		ctx:   ctx,
	}
}

type EventIns struct {
	Event
	ctx context.Context
}

func (e *EventIns) Context() context.Context {
	return e.ctx
}

type EventHandle struct {
	Event
	handle func(event IEventIns)
}

func (e *EventHandle) Handle(event IEventIns) {
	if e.handle != nil {
		e.handle(event)
	}
}

type EventPublisher struct {
	eventBus IEventBus
}

func (ep *EventPublisher) Init(bus IEventBus) {
	if ep.eventBus != nil || bus == nil {
		return
	}
	ep.eventBus = bus
}

func (ep *EventPublisher) PubEvent(events ...IEventIns) {
	if ep.eventBus == nil {
		return
	}
	ep.eventBus.Publish(events...)
}

type EventSubscriber struct {
	subscriber string
	eventBus   IEventBus
	handle     map[string]IEventHandle
}

func (es *EventSubscriber) Init(name string, bus IEventBus, events ...IEventHandle) {
	if es.subscriber != "" || name == "" {
		return
	}
	if es.eventBus != nil || bus == nil {
		return
	}
	es.handle = make(map[string]IEventHandle, len(events)*2)
	es.subscriber = name
	es.eventBus = bus
	es.Subscribe(events...)
}

func (es *EventSubscriber) Subscriber() string {
	return es.subscriber
}

func (es *EventSubscriber) Subscribe(events ...IEventHandle) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	for _, event := range events {
		es.eventBus.Register(es, event)
		es.handle[event.EventName()] = event
	}
}

func (es *EventSubscriber) UnSubscribe(events ...IEvent) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	es.eventBus.UnRegister(es, events...)
}

func (es *EventSubscriber) OnEvent(event IEventIns) {
	if es.handle == nil {
		return
	}
	handle := es.handle[event.EventName()]
	if handle != nil {
		handle.Handle(event)
	}
}

type EventSubscriberSingle struct {
	EventSubscriber
}

func NewEventSubscriberSingle(name string, bus IEventBus, events ...IEventHandle) IEventSubscriber {
	t := &EventSubscriberSingle{}
	t.Init(name, bus, events...)
	return t
}

type EventSubscriberSafety struct {
	EventSubscriber
	handleLock sync.RWMutex
}

func NewEventSubscriberSafety(name string, bus IEventBus, events ...IEventHandle) IEventSubscriber {
	e := &EventSubscriberSafety{}
	e.Init(name, bus, events...)
	return e
}

func (es *EventSubscriberSafety) Subscribe(events ...IEventHandle) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	es.handleLock.Lock()
	defer es.handleLock.Unlock()
	for _, event := range events {
		es.eventBus.Register(es, event)
		es.handle[event.EventName()] = event
	}
}

func (es *EventSubscriberSafety) UnSubscribe(events ...IEvent) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	es.handleLock.Lock()
	defer es.handleLock.Unlock()
	es.eventBus.UnRegister(es, events...)
	for _, event := range events {
		delete(es.handle, event.EventName())
	}
}

func (es *EventSubscriberSafety) OnEvent(event IEventIns) {
	if es.handle == nil {
		return
	}
	es.handleLock.RLock()
	defer es.handleLock.RUnlock()
	handle := es.handle[event.EventName()]
	if handle != nil {
		handle.Handle(event)
	}
}

type EventBus struct {
	subscribers map[string]map[string]IEventSubscriber // event name -> subscriber name -> subscriber
}

func (eb *EventBus) Register(sub IEventSubscriber, events ...IEvent) {
	for _, event := range events {
		subs := eb.subscribers[event.EventName()]
		if subs == nil {
			subs = make(map[string]IEventSubscriber)
		}
		subs[sub.Subscriber()] = sub
		eb.subscribers[event.EventName()] = subs
	}
}

func (eb *EventBus) UnRegister(sub IEventSubscriber, events ...IEvent) {
	for _, event := range events {
		subs := eb.subscribers[event.EventName()]
		if subs != nil {
			delete(subs, sub.Subscriber())
		}
	}
}

func (eb *EventBus) Publish(events ...IEventIns) {
	for _, event := range events {
		subs := eb.subscribers[event.EventName()]
		if subs != nil {
			for _, sub := range subs {
				sub.OnEvent(event)
			}
		}
	}
}

type EventBusSingle struct {
	EventBus
}

func NewEventBusSingle() IEventBus {
	return &EventBusSingle{
		EventBus: EventBus{
			subscribers: make(map[string]map[string]IEventSubscriber),
		},
	}
}

type eventBusWithLock struct {
	EventBus
	lock sync.RWMutex
}

func newEventBusWithLock() IEventBus {
	return &eventBusWithLock{
		EventBus: EventBus{
			subscribers: make(map[string]map[string]IEventSubscriber),
		},
	}
}

func (eb *eventBusWithLock) Register(sub IEventSubscriber, events ...IEvent) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.EventBus.Register(sub, events...)
}

func (eb *eventBusWithLock) UnRegister(sub IEventSubscriber, events ...IEvent) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.EventBus.UnRegister(sub, events...)
}

func (eb *eventBusWithLock) Publish(events ...IEventIns) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()
	eb.EventBus.Publish(events...)
}

type EventBusSafety struct {
	buckets []IEventBus
}

func NewEventBusBucket(bucketNum int) IEventBus {
	eb := &EventBusSafety{
		buckets: make([]IEventBus, bucketNum),
	}
	for i := 0; i < bucketNum; i++ {
		eb.buckets[i] = newEventBusWithLock()
	}
	return eb
}

func (eb *EventBusSafety) Register(sub IEventSubscriber, events ...IEvent) {
	eb.hashInvoke(func(bus IEventBus, event IEvent) {
		bus.Register(sub, event)
	}, events...)
}

func (eb *EventBusSafety) UnRegister(sub IEventSubscriber, events ...IEvent) {
	eb.hashInvoke(func(bus IEventBus, event IEvent) {
		bus.UnRegister(sub, event)
	}, events...)
}

func (eb *EventBusSafety) Publish(events ...IEventIns) {
	iEvents := make([]IEvent, len(events))
	for i, event := range events {
		iEvents[i] = event
	}
	eb.hashInvoke(func(bus IEventBus, event IEvent) {
		eventIns, ok := event.(IEventIns)
		if !ok {
			return
		}
		bus.Publish(eventIns)
	}, iEvents...)
}

func (eb *EventBusSafety) hashInvoke(f func(IEventBus, IEvent), events ...IEvent) {
	for _, event := range events {
		hash32a := fnv.New32a()
		n, err := hash32a.Write([]byte(event.EventName()))
		if err != nil {
			_ = n
			continue
		}
		hashValue := hash32a.Sum32()
		idx := int(hashValue) % len(eb.buckets)
		f(eb.buckets[idx], event)
	}
}
