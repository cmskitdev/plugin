package plugins

import (
	"testing"
	"time"
)

// MockPlugin implements the Plugin interface for testing
type MockPlugin struct {
	id       string
	handlers map[string]EventHandler[any]
	received []Event
}

func (m *MockPlugin) ID() string {
	return m.id
}

func (m *MockPlugin) Init() (PluginRegistration, error) {
	return PluginRegistration{
		Events: []Event{EventInit, EventShutdown},
	}, nil
}

func (m *MockPlugin) Receive(e Event) {
	m.received = append(m.received, e)
}

func (m *MockPlugin) Handlers() map[string]EventHandler[any] {
	return m.handlers
}

func NewMockPlugin(id string) *MockPlugin {
	return &MockPlugin{
		id: id,
		handlers: map[string]EventHandler[any]{
			string(EventInit): func(msg Message[any]) error {
				return nil
			},
			string(EventShutdown): func(msg Message[any]) error {
				return nil
			},
		},
		received: []Event{},
	}
}

func TestRegistryRegister(t *testing.T) {
	bus := NewEventBus[any]()
	registry := NewRegistry(bus)
	
	plugin := NewMockPlugin("test-plugin")
	
	// Register the plugin
	registry.Register(plugin)
	
	// Verify plugin is registered
	retrievedPlugin, exists := registry.Get("test-plugin")
	if !exists {
		t.Fatal("Plugin was not registered")
	}
	
	if retrievedPlugin.ID() != "test-plugin" {
		t.Errorf("Expected plugin ID 'test-plugin', got '%s'", retrievedPlugin.ID())
	}
}

func TestRegistryUnregister(t *testing.T) {
	bus := NewEventBus[any]()
	registry := NewRegistry(bus)
	
	plugin := NewMockPlugin("test-plugin")
	
	// Register then unregister the plugin
	registry.Register(plugin)
	registry.Unregister("test-plugin")
	
	// Verify plugin is unregistered
	_, exists := registry.Get("test-plugin")
	if exists {
		t.Fatal("Plugin should have been unregistered")
	}
}

func TestEventBusSubscribeAndPublish(t *testing.T) {
	bus := NewEventBus[any]()
	
	// Create a channel to capture handler execution
	handlerCalled := make(chan bool, 1)
	
	handler := func(msg Message[any]) error {
		handlerCalled <- true
		return nil
	}
	
	// Subscribe to the event
	bus.Subscribe(EventInit, "test-plugin", handler)
	
	// Publish the event
	bus.Publish(Message[any]{
		Event: EventInit,
		Data:  "test data",
	})
	
	// Wait for handler to be called
	select {
	case <-handlerCalled:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Handler was not called within timeout")
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	bus := NewEventBus[any]()
	
	handlerCalled := make(chan bool, 1)
	
	handler := func(msg Message[any]) error {
		handlerCalled <- true
		return nil
	}
	
	// Subscribe then unsubscribe
	bus.Subscribe(EventInit, "test-plugin", handler)
	bus.Unsubscribe(EventInit, "test-plugin")
	
	// Publish the event
	bus.Publish(Message[any]{
		Event: EventInit,
		Data:  "test data",
	})
	
	// Handler should not be called
	select {
	case <-handlerCalled:
		t.Fatal("Handler should not have been called after unsubscribe")
	case <-time.After(50 * time.Millisecond):
		// Success - handler was not called
	}
}

func TestPluginRegistrationWithEventBus(t *testing.T) {
	bus := NewEventBus[any]()
	registry := NewRegistry(bus)
	
	// Create channels to capture handler execution
	initHandlerCalled := make(chan bool, 1)
	shutdownHandlerCalled := make(chan bool, 1)
	
	plugin := &MockPlugin{
		id: "test-plugin",
		handlers: map[string]EventHandler[any]{
			string(EventInit): func(msg Message[any]) error {
				initHandlerCalled <- true
				return nil
			},
			string(EventShutdown): func(msg Message[any]) error {
				shutdownHandlerCalled <- true
				return nil
			},
		},
		received: []Event{},
	}
	
	// Register the plugin
	registry.Register(plugin)
	
	// Publish init event
	bus.Publish(Message[any]{
		Event: EventInit,
		Data:  "init data",
	})
	
	// Verify init handler was called
	select {
	case <-initHandlerCalled:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Init handler was not called within timeout")
	}
	
	// Publish shutdown event
	bus.Publish(Message[any]{
		Event: EventShutdown,
		Data:  "shutdown data",
	})
	
	// Verify shutdown handler was called
	select {
	case <-shutdownHandlerCalled:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Shutdown handler was not called within timeout")
	}
}

func TestMultiplePluginsEventHandling(t *testing.T) {
	bus := NewEventBus[any]()
	registry := NewRegistry(bus)
	
	// Create channels to capture handler execution
	plugin1Called := make(chan bool, 1)
	plugin2Called := make(chan bool, 1)
	
	plugin1 := &MockPlugin{
		id: "plugin1",
		handlers: map[string]EventHandler[any]{
			string(EventInit): func(msg Message[any]) error {
				plugin1Called <- true
				return nil
			},
		},
		received: []Event{},
	}
	
	plugin2 := &MockPlugin{
		id: "plugin2",
		handlers: map[string]EventHandler[any]{
			string(EventInit): func(msg Message[any]) error {
				plugin2Called <- true
				return nil
			},
		},
		received: []Event{},
	}
	
	// Register both plugins
	registry.Register(plugin1)
	registry.Register(plugin2)
	
	// Publish init event
	bus.Publish(Message[any]{
		Event: EventInit,
		Data:  "init data",
	})
	
	// Verify both handlers were called
	select {
	case <-plugin1Called:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Plugin1 handler was not called within timeout")
	}
	
	select {
	case <-plugin2Called:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Plugin2 handler was not called within timeout")
	}
}