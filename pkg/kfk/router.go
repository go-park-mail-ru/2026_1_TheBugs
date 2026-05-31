package kfk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Router struct {
	handlers map[string]Handler
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]Handler)}
}

func (r *Router) Register(eventType string, h Handler) {
	r.handlers[eventType] = h
}

func (r *Router) Handle(ctx context.Context, msg []byte) error {
	var e Event
	if err := json.Unmarshal(msg, &e); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	h, ok := r.handlers[e.Type]
	if !ok {
		return fmt.Errorf("no handler for event type: %s", e.Type)
	}

	return h(ctx, e.Payload)
}
