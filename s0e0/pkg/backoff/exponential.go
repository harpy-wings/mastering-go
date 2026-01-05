package backoff

import (
	"context"
	"fmt"
	"math"
	"time"
)

type config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func ExponentialDo(ctx context.Context, fn func(ctx context.Context) error, opts ...Option) error {
	config := &config{
		MaxAttempts:  5,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     10 * time.Second,
	}
	for _, opt := range opts {
		opt(config)
	}
	for i := 0; i < config.MaxAttempts; i++ {
		if err := fn(ctx); err == nil {
			return nil
		}
		time.Sleep(time.Duration(math.Pow(float64(i+1), 3)) * config.InitialDelay)
	}
	return fmt.Errorf("failed to execute function after %d attempts", config.MaxAttempts)
}

func ExponentialDoWithReturn[T any](ctx context.Context, fn func(ctx context.Context) (T, error), opts ...Option) (T, error) {
	config := &config{
		MaxAttempts:  5,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     10 * time.Second,
	}
	for _, opt := range opts {
		opt(config)
	}
	for i := 0; i < config.MaxAttempts; i++ {
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}
		time.Sleep(time.Duration(math.Pow(float64(i+1), 3)) * config.InitialDelay)
	}
	return *new(T), fmt.Errorf("failed to execute function after %d attempts", config.MaxAttempts)
}
