package plugins

import (
	"sync"
)

// Event is a type of event that can be published to the event bus.
type Event string

const (
	// EventInit is the event that is published when the plugin is initialized.
	EventInit Event = "init"
	// EventShutdown is the event that is published when the plugin is shutdown.
	EventShutdown Event = "shutdown"
	// EventDispatch is the event that is published when the plugin is dispatched.
	EventDispatch Event = "dispatch"
)

// Message is a message that can be published to the event bus.
type Message[T any] struct {
	Event Event
	Data  T
}

// EventHandler is a function that handles a message.
type EventHandler[T any] func(Message[T]) error

// EventBus is a bus for publishing and subscribing to events.
type EventBus[T any] struct {
	mu       sync.RWMutex
	handlers map[Event]map[string]EventHandler[T] // eventName → pluginID → handler
}

// NewEventBus creates a new event bus.
//
// Arguments:
//   - T: the type of the message
//
// Returns:
//   - *EventBus[T]: the event bus
func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{handlers: make(map[Event]map[string]EventHandler[T])}
}

// Subscribe subscribes to an event.
func (b *EventBus[T]) Subscribe(e Event, pluginID string, handler EventHandler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.handlers[e] == nil {
		b.handlers[e] = make(map[string]EventHandler[T])
	}
	b.handlers[e][pluginID] = handler
}

// Unsubscribe unsubscribes from an event.
func (b *EventBus[T]) Unsubscribe(e Event, pluginID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if h := b.handlers[e]; h != nil {
		delete(h, pluginID)
	}
}

// Publish publishes a message to the event bus.
func (b *EventBus[T]) Publish(e Message[T]) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers[e.Event] {
		go func(h EventHandler[T]) {
			_ = h(e) // optionally log errors
		}(h)
	}
}
