package main

import (
	"context"
	"fmt"
)

// Handler defines the abstract interface for handlers
type Handler interface {
	Handle(ctx context.Context, request string) (string, error)
}

// ConcreteHandlerA implements the Handler interface
type ConcreteHandlerA struct{}

func (h *ConcreteHandlerA) Handle(ctx context.Context, request string) (string, error) {
	if request == "A" {
		return "Handled by ConcreteHandlerA", nil
	}
	return "", fmt.Errorf("ConcreteHandlerA cannot handle the request")
}

// ConcreteHandlerB implements the Handler interface
type ConcreteHandlerB struct{}

func (h *ConcreteHandlerB) Handle(ctx context.Context, request string) (string, error) {
	if request == "B" {
		return "Handled by ConcreteHandlerB", nil
	}
	return "", fmt.Errorf("ConcreteHandlerB cannot handle the request")
}

func main() {
	handlerA := &ConcreteHandlerA{}
	handlerB := &ConcreteHandlerB{}

	handlerA.Handle(context.Background(), "A")
	handlerB.Handle(context.Background(), "B")
}
