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

func NewEventPublisher(bus IEventBus) IEventPublisher {
	t := &EventPublisher{}
	t.Init(bus)
	return t
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

type eventSubscriber struct {
	subscriber string
	eventBus   IEventBus
	handle     map[string]IEventHandle
}

func (es *eventSubscriber) Init(name string, bus IEventBus, events ...IEventHandle) {
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

func (es *eventSubscriber) Subscriber() string {
	return es.subscriber
}

func (es *eventSubscriber) Subscribe(events ...IEventHandle) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	for _, event := range events {
		es.eventBus.Register(es, event)
		es.handle[event.EventName()] = event
	}
}

func (es *eventSubscriber) UnSubscribe(events ...IEvent) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	es.eventBus.UnRegister(es, events...)
}

func (es *eventSubscriber) OnEvent(event IEventIns) {
	if es.handle == nil {
		return
	}
	handle := es.handle[event.EventName()]
	if handle != nil {
		handle.Handle(event)
	}
}

type EventSubscriberSingle struct {
	eventSubscriber
}

func NewEventSubscriberSingle(name string, bus IEventBus, events ...IEventHandle) IEventSubscriber {
	e := &EventSubscriberSingle{}
	e.Init(name, bus, events...)
	return e
}

type EventSubscriberSafety struct {
	eventSubscriber
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
	var iEvents = make([]IEvent, 0, len(events))
	func() {
		es.handleLock.Lock()
		defer es.handleLock.Unlock()
		for _, event := range events {
			es.handle[event.EventName()] = event
			iEvents = append(iEvents, event)
		}
	}()
	es.eventBus.Register(es, iEvents...)
}

func (es *EventSubscriberSafety) UnSubscribe(events ...IEvent) {
	if es.eventBus == nil || len(events) == 0 {
		return
	}
	es.eventBus.UnRegister(es, events...)
	func() {
		es.handleLock.Lock()
		defer es.handleLock.Unlock()
		for _, event := range events {
			delete(es.handle, event.EventName())
		}
	}()
}

func (es *EventSubscriberSafety) OnEvent(event IEventIns) {
	if es.handle == nil {
		return
	}
	func() {
		es.handleLock.RLock()
		defer es.handleLock.RUnlock()
		handle := es.handle[event.EventName()]
		if handle != nil {
			handle.Handle(event)
		}
	}()
}

type eventBus struct {
	subscribers map[string]map[string]IEventSubscriber // event name -> subscriber name -> subscriber
}

func newEventBus() IEventBus {
	return &eventBus{
		subscribers: make(map[string]map[string]IEventSubscriber, 128),
	}
}

func (eb *eventBus) Register(sub IEventSubscriber, events ...IEvent) {
	for _, event := range events {
		subs := eb.subscribers[event.EventName()]
		if subs == nil {
			subs = make(map[string]IEventSubscriber, 128)
		}
		subs[sub.Subscriber()] = sub
		eb.subscribers[event.EventName()] = subs
	}
}

func (eb *eventBus) UnRegister(sub IEventSubscriber, events ...IEvent) {
	for _, event := range events {
		subs := eb.subscribers[event.EventName()]
		if subs != nil {
			delete(subs, sub.Subscriber())
		}
	}
}

func (eb *eventBus) Publish(events ...IEventIns) {
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
	IEventBus
}

func NewEventBusSingle() IEventBus {
	return &EventBusSingle{
		IEventBus: newEventBus(),
	}
}

type eventBusWithLock struct {
	IEventBus
	lock sync.RWMutex
}

func newEventBusWithLock() IEventBus {
	return &eventBusWithLock{
		IEventBus: newEventBus(),
	}
}

func (eb *eventBusWithLock) Register(sub IEventSubscriber, events ...IEvent) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.IEventBus.Register(sub, events...)
}

func (eb *eventBusWithLock) UnRegister(sub IEventSubscriber, events ...IEvent) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.IEventBus.UnRegister(sub, events...)
}

func (eb *eventBusWithLock) Publish(events ...IEventIns) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()
	eb.IEventBus.Publish(events...)
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
