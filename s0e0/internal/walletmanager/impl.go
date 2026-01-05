package walletmanager

import (
	"context"
	"time"

	"github.com/harpy-wings/mastering-go/s0e0/internal/usermanager"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type walletManager struct {
	logger logrus.FieldLogger
	// tracer trace.Tracer
	rdb redis.UniversalClient
	um  usermanager.IUserManager

	config struct {
		ttl time.Duration
	}
}

var _ IWalletManager = (*walletManager)(nil)

func New(ctx context.Context, opts ...Option) (IWalletManager, error) {
	wm := new(walletManager)
	wm.setDefaults(ctx)
	for _, opt := range opts {
		if err := opt(wm); err != nil {
			return nil, err
		}
	}
	err := wm.init(ctx)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

func (wm *walletManager) setDefaults(ctx context.Context) {
	wm.config.ttl = 1 * time.Hour
}

func (wm *walletManager) init(ctx context.Context) error {
	return nil
}
