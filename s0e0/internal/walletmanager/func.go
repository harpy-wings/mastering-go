package walletmanager

import "context"

// GetBalanceOf gets the balance of a user
// this is a frequent operation so it should be focused on performance.
// ferquently called allowed.
func (wm *walletManager) GetBalanceOf(ctx context.Context, uuid string) (float64, error) {
	return 0, nil
}

// Transfer transfers money between two users
// this is critical flow of our application, since it is ferquently called, and transactions should be atomic with fully data consistency.
func (wm *walletManager) Transfer(ctx context.Context, fromUUID string, toUUID string, amount float64) error {
	return nil
}

func (wm *walletManager) GracefulStop(ctx context.Context) error {
	wm.logger.Info("stopping wallet manager")
	return nil
}
