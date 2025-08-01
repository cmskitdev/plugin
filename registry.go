package plugin

import (
	"context"
	"sync"

	"github.com/cmskitdev/engine"
)

type PluginConfig struct {
	Name           string   `json:"name"`
	Version        string   `json:"version"`         // Semantic version
	CompatibleWith []string `json:"compatible_with"` // e.g. ["v1", "v2"]
	Description    string   `json:"description"`
	Author         string   `json:"author,omitempty"`
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]PluginConfig),
	}
}

type PluginRegistry struct {
	plugins map[string]PluginConfig
	mu      sync.RWMutex
}

func (r *PluginRegistry) Register(name string, plugin PluginConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[name] = plugin
}

type Plugin interface {
	Transform(ctx context.Context, item engine.DataItemContainer) (engine.DataItemContainer, error)
	GetConfig() PluginConfig
}
