package kfk

import (
	"context"
)

type Event struct {
	Type    string `json:"type"`
	Version string `json:"version,omitempty"`
	Payload any    `json:"payload"`
}

type Handler func(context.Context, any) error
