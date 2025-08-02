package plugins

import (
	"fmt"
	"sync"
)

type Registry[T any] struct {
	mu      sync.RWMutex
	plugins map[string]Plugin[T]
	bus     *EventBus[T]
}

func NewRegistry[T any](bus *EventBus[T]) *Registry[T] {
	return &Registry[T]{plugins: map[string]Plugin[T]{}, bus: bus}
}

func (r *Registry[T]) Register(p Plugin[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[p.ID()] = p
	
	// Register plugin handlers with the event bus
	handlers := p.Handlers()
	for eventName, handler := range handlers {
		r.bus.Subscribe(Event(eventName), p.ID(), handler)
		fmt.Printf("Subscribed plugin %s to event %s\n", p.ID(), eventName)
	}
	
	fmt.Printf("Registered plugin %s, plugins: %d\n", p.ID(), len(r.plugins))
}

func (r *Registry[T]) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.plugins[id]; ok {
		delete(r.plugins, id)
	}
}

func (r *Registry[T]) Get(id string) (Plugin[T], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	return p, ok
}
