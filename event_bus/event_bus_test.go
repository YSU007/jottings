package event_bus

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type TestPublisher struct {
	EventPublisher
}

func NewTestPublisher(bus IEventBus) *TestPublisher {
	t := &TestPublisher{}
	t.Init(bus)
	return t
}

func EventHandle1(event IEventIns) {
	fmt.Println("receive event", event.EventName(), event.Context().Value("testK"))
}

func TestEventBusSingle(t *testing.T) {
	testEvent := Event{name: "TestEvent"}

	bus := NewEventBus(BusSingle)
	p := NewTestPublisher(bus)
	s := NewEventSubscriber(SubscriberSingle, "TestSubscriber", bus, &EventHandle{testEvent, EventHandle1})
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV1")))
	s.UnSubscribe(&testEvent)
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
	s.Subscribe(&EventHandle{testEvent, EventHandle1})
	p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
}

func TestEventBusSafety(t *testing.T) {
	bus := NewEventBus(BusSafety)
	s := NewEventSubscriber(SubscriberSafety, "TestSubscriber", bus)
	for i := 0; i < 100; i++ {
		j := i
		testEvent := Event{name: fmt.Sprint("TestEvent", j)}
		s.Subscribe(&EventHandle{testEvent, EventHandle1})
		go func() {
			p := NewTestPublisher(bus)
			p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV1")))
			p.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
		}()
	}
	time.Sleep(time.Second)
}
