package backoff

import "time"

type Option func(*config)

func WithMaxAttempts(maxAttempts int) Option {
	return func(c *config) {
		c.MaxAttempts = maxAttempts
	}
}

func WithInitialDelay(initialDelay time.Duration) Option {
	return func(c *config) {
		c.InitialDelay = initialDelay
	}
}

func WithMaxDelay(maxDelay time.Duration) Option {
	return func(c *config) {
		c.MaxDelay = maxDelay
	}
}
