package walletmanager

import (
	"context"
	"time"

	"github.com/harpy-wings/mastering-go/s0e0/internal/usermanager"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Option func(*walletManager) error

func WithLogger(logger logrus.FieldLogger) Option {
	return func(wm *walletManager) error {
		wm.logger = logger
		return nil
	}
}

func WithRedis(rdb redis.UniversalClient) Option {
	return func(wm *walletManager) error {
		err := rdb.Ping(context.Background()).Err()
		if err != nil {
			return err
		}
		wm.rdb = rdb
		return nil
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(wm *walletManager) error {
		wm.config.ttl = ttl
		return nil
	}
}

func WithUserManager(um usermanager.IUserManager) Option {
	return func(wm *walletManager) error {
		wm.um = um
		return nil
	}
}
