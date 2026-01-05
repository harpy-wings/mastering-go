package usermanager

import (
	"context"

	"github.com/harpy-wings/mastering-go/s0e0/pkg/models"
)

type IUserManager interface {
	// ListUsers lists all users with pagination
	// only admin call this method so it doesn't need focuse of performance.
	ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error)
	// CreateUser creates a new user
	// onboarding new users is not a frequent operation so it doesn't need focus of performance. instead, it should be focused on validation and data consistency.
	CreateUser(ctx context.Context, user *models.User) error
	// GetUserByUUID gets a user by their UUID
	// this is a frequent operation so it should be focused on performance.
	// ferquently called allowed.
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	// GetUserByEmail gets a user by their email
	// this is a frequent operation so it should be focused on performance.
	// ferquently called allowed.
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// UpdateUser updates a user
	// this is not a frequent operation so it doesn't need focus of performance. instead, it should be focused on validation and data consistency.
	UpdateUser(ctx context.Context, user *models.User) error
	// DeleteUser deletes a user
	// this is not a frequent operation so it doesn't need focus of performance. instead, it should be focused on validation and data consistency.
	DeleteUser(ctx context.Context, uuid string) error
}

//go:generate mockgen -source=interface.go -destination=./../../tests/mocks/user_manager_mock.go -package=mocks
