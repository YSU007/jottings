package event_bus

import (
	"context"
	"fmt"
	"testing"
)

type TestPublisher struct {
	EventPublisher
}

func NewTestPublisher(bus IEventBus) *TestPublisher {
	t := &TestPublisher{}
	t.Init(bus)
	return t
}

type TestSubscriber struct {
	EventSubscriber
}

func NewTestSubscriber(name string, bus IEventBus, events ...IEventHandle) *TestSubscriber {
	t := &TestSubscriber{}
	t.Init(name, bus, events...)
	return t
}

func EventHandle1(event IEventIns) {
	fmt.Println("receive event", event.EventName(), event.Context().Value("testK"))
}

func TestEventBusSingle(t *testing.T) {
	testEvent := Event{name: "TestEvent"}

	bus := NewEventBus(Single)
	p := NewTestPublisher(bus)
	s := NewTestSubscriber("TestSubscriber", bus, &EventHandle{testEvent, EventHandle1})
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV1")))
	s.UnSubscribe(&testEvent)
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
	s.Subscribe(&EventHandle{testEvent, EventHandle1})
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
}
