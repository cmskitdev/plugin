package plugins

import "sync"

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
	p.Subscribe(r.bus)
}

func (r *Registry[T]) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.plugins[id]; ok {
		p.Unsubscribe(r.bus)
		delete(r.plugins, id)
	}
}

func (r *Registry[T]) Get(id string) (Plugin[T], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	return p, ok
}
