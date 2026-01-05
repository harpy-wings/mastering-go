package walletmanager

import "context"

// IWalletManager is the interface for the wallet manager.
type IWalletManager interface {
	// GetBalanceOf gets the balance of a user
	// this is a frequent operation so it should be focused on performance.
	// ferquently called allowed.
	GetBalanceOf(ctx context.Context, uuid string) (float64, error)
	// Transfer transfers money between two users
	// this is critical flow of our application, since it is ferquently called, and transactions should be atomic with fully data consistency.
	Transfer(ctx context.Context, fromUUID string, toUUID string, amount float64) error

	GracefulStop(ctx context.Context) error
}

//go:generate mockgen -source=interface.go -destination=./../../tests/mocks/wallet_manager_mock.go -package=mocks
