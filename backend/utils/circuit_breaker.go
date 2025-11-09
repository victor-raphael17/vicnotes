package utils

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// StateClosed means the circuit is working normally
	StateClosed CircuitState = iota
	// StateOpen means the circuit has failed and is rejecting requests
	StateOpen
	// StateHalfOpen means the circuit is testing if it can recover
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	failureThreshold int
	successThreshold int
	timeout         time.Duration
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

// Call executes a function through the circuit breaker
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if we should transition from Open to HalfOpen
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// Execute the function
	err := fn()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Transition to Open if threshold exceeded
		if cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		}

		return err
	}

	// Success
	cb.failureCount = 0

	if cb.state == StateHalfOpen {
		cb.successCount++
		// Transition back to Closed if success threshold met
		if cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
		}
	}

	return nil
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
}

// StateString returns a string representation of the circuit state
func (state CircuitState) String() string {
	switch state {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}
