package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay time.Duration
	Multiplier float64
}

// DefaultRetryConfig returns sensible defaults for retries
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay: 5 * time.Second,
		Multiplier: 2.0,
	}
}

// RetryFunc is a function that can be retried
type RetryFunc func() error

// FallbackFunc is a fallback function to execute if retries fail
type FallbackFunc func(lastErr error) error

// Retry executes a function with exponential backoff retry logic
func Retry(fn RetryFunc, config RetryConfig) error {
	var lastErr error
	
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Don't sleep after the last failed attempt
		if attempt < config.MaxAttempts-1 {
			delay := calculateBackoff(attempt, config)
			time.Sleep(delay)
		}
	}
	
	return fmt.Errorf("retry failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// RetryWithFallback executes a function with retries and falls back if all retries fail
func RetryWithFallback(fn RetryFunc, fallback FallbackFunc, config RetryConfig) error {
	var lastErr error
	
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Don't sleep after the last failed attempt
		if attempt < config.MaxAttempts-1 {
			delay := calculateBackoff(attempt, config)
			time.Sleep(delay)
		}
	}
	
	// All retries failed, try fallback
	if fallback != nil {
		return fallback(lastErr)
	}
	
	return fmt.Errorf("retry failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// calculateBackoff calculates exponential backoff with jitter
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	// Exponential backoff: initialDelay * (multiplier ^ attempt)
	backoff := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))
	
	// Cap at max delay
	if backoff > float64(config.MaxDelay) {
		backoff = float64(config.MaxDelay)
	}
	
	// Add jitter (±10%)
	jitter := backoff * 0.1 * (2*rand.Float64() - 1)
	
	return time.Duration(backoff + jitter)
}
