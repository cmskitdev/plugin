// Package plugin provides provides the ability for registering plugins.
package plugin

import (
	"context"

	"github.com/cmskitdev/common"
)

// Plugin is the interface that all plugins must implement.
type Plugin interface {
	Transform(ctx context.Context, item interface{}) (interface{}, error)
	Export(ctx context.Context, objectType common.ObjectType, object interface{}) (interface{}, error)
	GetConfig() PluginConfig
}

// Config is the configuration for a plugin instance.
type Config struct {
	ID             string
	EnableReporter bool
	Reporter       *Reporter
}

// Base is the instance of a plugin.
type Base struct {
	Config   Config
	Reporter *Reporter
}

// NewPlugin creates a new plugin instance.
//
// Arguments:
// - config: The plugin configuration.
//
// Returns:
// - The plugin instance.
func NewPlugin(config Config) *Base {
	var reporter *Reporter
	if config.EnableReporter {
		reporter = NewReporter(nil, config.Reporter.Interval, config.Reporter.BatchSize)
	}

	return &Base{
		Config:   config,
		Reporter: reporter,
	}
}
