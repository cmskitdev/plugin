package plugin

import (
	"sync"
)

type PluginConfig struct {
	Name           string   `json:"name"`
	Version        string   `json:"version"`         // Semantic version
	CompatibleWith []string `json:"compatible_with"` // e.g. ["v1", "v2"]
	Description    string   `json:"description"`
	Author         string   `json:"author,omitempty"`
}

func NewPluginRegistry() *PluginRegistry {
	registry := &PluginRegistry{
		plugins: make(map[string]interface{}),
	}

	registry.Register("redis", &redis.Plugin{})

	return registry
}

type PluginRegistry struct {
	plugins map[string]interface{}
	mu      sync.RWMutex
}

func (r *PluginRegistry) Register(name string, plugin interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[name] = plugin
}
