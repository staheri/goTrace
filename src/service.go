package main

import (
	"context"
)

// Service is a service that depends on the Foo service
type Service struct {}

// Foo is an external service
type Foo interface {
	// Bar is an RPC that can be cancelled via it's context
	Bar(ctx context.Context, arg string)
}

// Baz is an endpoint on our service
func (*Service) Baz(ctx context.Context, mock Foo) {
	done := make(chan struct{})
	go func() {
                // This might result in a runtime.Goexit() which means the message on the done chan is never sent
		mock.Bar(ctx, "test")
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}
