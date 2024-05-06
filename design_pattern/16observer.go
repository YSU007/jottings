package main

import (
	"fmt"
	"sync"
)

// Observer 定义观察者接口
type Observer interface {
	Update(subject Subject)
}

// Subject 定义主题接口
type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Notify()
}

// ConcreteObserver 实现观察者
type ConcreteObserver struct {
	name string
}

func (o *ConcreteObserver) Update(subject Subject) {
	fmt.Printf("Observer %s received an update from the subject\n", o.name)
}

// ConcreteSubject 实现主题
type ConcreteSubject struct {
	observers map[Observer]struct{}
	mux       sync.RWMutex
}

func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{
		observers: make(map[Observer]struct{}),
	}
}

func (s *ConcreteSubject) Register(observer Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.observers[observer] = struct{}{}
}

func (s *ConcreteSubject) Unregister(observer Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.observers, observer)
}

func (s *ConcreteSubject) Notify() {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for observer := range s.observers {
		observer.Update(s)
	}
}

func main() {
	subject := NewConcreteSubject()

	observerA := &ConcreteObserver{name: "ObserverA"}
	observerB := &ConcreteObserver{name: "ObserverB"}

	subject.Register(observerA)
	subject.Register(observerB)

	subject.Notify()
}
