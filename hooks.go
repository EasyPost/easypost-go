package easypost

import (
	"context"
	"github.com/google/uuid"
	"io"
	"net/url"
	"time"
)

// HookEvent is the base type for all hook events
type HookEvent struct {
}

// HookEventSubscriber is the base type for all hook event subscribers.
// ID needs to be unique in order to find and remove the subscriber
type HookEventSubscriber struct {
	ID string // need unique ID for removal
}

// RequestHookEvent is the data passed to a RequestHookEventSubscriberCallback function
type RequestHookEvent struct {
	HookEvent        // implements HookEvent
	Method           string
	Url              *url.URL
	RequestBody      io.ReadCloser
	Headers          map[string][]string
	RequestTimestamp time.Time
	Id               uuid.UUID
}

// RequestHookEventSubscriberCallback is the type of the callback function executed by an RequestHookEventSubscriber
type RequestHookEventSubscriberCallback func(ctx context.Context, event RequestHookEvent) error

// RequestHookEventSubscriber is a HookEventSubscriber that executes a RequestHookEventSubscriberCallback function when a RequestHookEvent is fired
type RequestHookEventSubscriber struct {
	HookEventSubscriber // implements HookEventSubscriber
	Callback            RequestHookEventSubscriberCallback
}

// Execute executes the RequestHookEventSubscriberCallback function of the RequestHookEventSubscriber
func (s RequestHookEventSubscriber) Execute(ctx context.Context, event RequestHookEvent) {
	err := s.Callback(ctx, event)
	if err != nil {
		return
	}
}

// ResponseHookEvent is the data passed to a ResponseHookEventSubscriberCallback function
type ResponseHookEvent struct {
	HookEvent         // implements HookEvent
	HttpStatus        int
	Method            string
	Url               *url.URL
	Headers           map[string][]string
	ResponseBody      io.ReadCloser
	RequestTimestamp  time.Time
	ResponseTimestamp time.Time
	Id                uuid.UUID
}

// ResponseHookEventSubscriberCallback is the type of the callback function executed by an ResponseHookEventSubscriber
type ResponseHookEventSubscriberCallback func(ctx context.Context, event ResponseHookEvent) error

// ResponseHookEventSubscriber is a HookEventSubscriber that executes a ResponseHookEventSubscriberCallback function when a ResponseHookEvent is fired
type ResponseHookEventSubscriber struct {
	HookEventSubscriber // implements HookEventSubscriber
	Callback            ResponseHookEventSubscriberCallback
}

// Execute executes the ResponseHookEventSubscriberCallback function of the ResponseHookEventSubscriber
func (s ResponseHookEventSubscriber) Execute(ctx context.Context, event ResponseHookEvent) {
	err := s.Callback(ctx, event)
	if err != nil {
		return
	}
}

// Hooks is a collection of HookEventSubscriber instances for various hooks available in the client
type Hooks struct {
	RequestHookEventSubscriptions  []RequestHookEventSubscriber // these are directly accessible by the user, but should not be modified directly (use Add/Remove methods)
	ResponseHookEventSubscriptions []ResponseHookEventSubscriber
}

// AddRequestEventSubscriber adds a RequestHookEventSubscriber to the Hooks instance to be executed when a RequestHookEvent is fired
func (h *Hooks) AddRequestEventSubscriber(subscriber RequestHookEventSubscriber) {
	h.RequestHookEventSubscriptions = append(h.RequestHookEventSubscriptions, subscriber)
}

// AddResponseEventSubscriber adds a ResponseHookEventSubscriber to the Hooks instance to be executed when a ResponseHookEvent is fired
func (h *Hooks) AddResponseEventSubscriber(subscriber ResponseHookEventSubscriber) {
	h.ResponseHookEventSubscriptions = append(h.ResponseHookEventSubscriptions, subscriber)
}

// RemoveRequestEventSubscriber removes a RequestHookEventSubscriber from the Hooks instance
func (h *Hooks) RemoveRequestEventSubscriber(subscriber RequestHookEventSubscriber) {
	for i, sub := range h.RequestHookEventSubscriptions {
		if sub.ID == subscriber.ID {
			h.RequestHookEventSubscriptions = append(h.RequestHookEventSubscriptions[:i], h.RequestHookEventSubscriptions[i+1:]...)
			return
		}
	}
}

// RemoveResponseEventSubscriber removes a ResponseHookEventSubscriber from the Hooks instance
func (h *Hooks) RemoveResponseEventSubscriber(subscriber ResponseHookEventSubscriber) {
	for i, sub := range h.ResponseHookEventSubscriptions {
		if sub.ID == subscriber.ID {
			h.ResponseHookEventSubscriptions = append(h.ResponseHookEventSubscriptions[:i], h.ResponseHookEventSubscriptions[i+1:]...)
			return
		}
	}
}
