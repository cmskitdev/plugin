package plugins

import (
	"sync"
	"time"
)

type State string

const (
	StateInit     State = "init"
	StateRunning  State = "running"
	StateShutdown State = "shutdown"
)

type StateTransition struct {
	State State
	Time  time.Time
}

type StateQueue struct {
	mu      sync.RWMutex
	history []StateTransition
}

func NewStateQueue() *StateQueue {
	return &StateQueue{history: []StateTransition{}}
}

func (s *StateQueue) Push(state State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, StateTransition{State: state, Time: time.Now()})
}

func (s *StateQueue) Pop() StateTransition {
	s.mu.Lock()
	defer s.mu.Unlock()
	transition := s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]
	return transition
}

func (s *StateQueue) Peek() StateTransition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.history[len(s.history)-1]
}

func (s *StateQueue) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.history)
}

func (s *StateQueue) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = []StateTransition{}
}

func (s *StateQueue) GetHistory() []StateTransition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.history
}
