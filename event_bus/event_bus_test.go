package event_bus

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var counter uint32

func incrementCounter(counter *uint32) {
	atomic.AddUint32(counter, 1)
}

func EventHandle1(event IEventIns) {
	incrementCounter(&counter)
	fmt.Println("receive event", event.EventName(), "count", counter, "value", event.Context().Value("testK"))
}

func TestEventBusSingle(t *testing.T) {
	testEvent := Event{name: "TestEvent"}

	bus := NewEventBus(BusSingle)
	pub := NewEventPublisher(bus)
	sub := NewEventSubscriber(SubscriberSingle, "TestSubscriber", bus, &EventHandle{testEvent, EventHandle1})
	pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV1")))
	sub.UnSubscribe(&testEvent)
	pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
	sub.Subscribe(&EventHandle{testEvent, EventHandle1})
	pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
}

func TestEventBusSafety(t *testing.T) {
	for i := 0; i < 10; i++ {
		bus := NewEventBus(BusSafety)
		sub := NewEventSubscriber(SubscriberSafety, "TestSubscriber", bus)
		for i := 0; i < 100; i++ {
			j := i
			testEvent := Event{name: fmt.Sprint("TestEvent", j)}
			sub.Subscribe(&EventHandle{testEvent, EventHandle1})
			go func() {
				pub := NewEventPublisher(bus)
				pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV1")))
				// sub.UnSubscribe(&testEvent)
				// pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
				// sub.Subscribe(&EventHandle{testEvent, EventHandle1})
				// pub.PubEvent(NewEventIns(testEvent, context.WithValue(context.Background(), "testK", "testV2")))
			}()
		}
		time.Sleep(time.Second)
	}
}
