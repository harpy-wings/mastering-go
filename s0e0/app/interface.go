package app

import "context"

// IApp is the interface for the application.
type IApp interface {
	// Run is the main function for the application.
	// It will lock the main thread and wait until the GracefullStop is called.
	Run(ctx context.Context) error

	// GracefulStop is the function that will be called to shutdown the application.
	// It will stop the application and wait for the main thread to exit.
	GracefulStop(ctx context.Context) error

	// HealthCheck is the function that will be called to check the health of the application.
	// It will return an error if the application is not healthy.
	// HealthCheck() error
}
