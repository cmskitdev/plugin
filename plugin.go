// Package plugins provides a framework for creating plugins.
package plugins

// Plugin is the interface that all plugins must implement.
type Plugin[T any] interface {
	ID() string
	Init() error
	Receive(e Event)
	Handlers() map[string]EventHandler[T]
}
