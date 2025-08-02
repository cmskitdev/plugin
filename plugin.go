// Package plugins provides a framework for creating plugins.
package plugins

// Plugin is the interface that all plugins must implement.
type Plugin[T any] interface {
	ID() string
	Init() (PluginRegistration, error)
	Receive(e Event)
	Handlers() map[Event]EventHandler[T]
}

type PluginRegistration struct {
	Events []Event
}
