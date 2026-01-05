package hermesapp

import (
	"context"

	"github.com/harpy-wings/mastering-go/s0e0/app"
)

type hermesApp struct {
}

var _ app.IApp = &hermesApp{}

func (a *hermesApp) Run(ctx context.Context) error {
	return nil
}

func (a *hermesApp) GracefulStop(ctx context.Context) error {
	return nil
}

func New(ctx context.Context) (app.IApp, error) {
	return &hermesApp{}, nil
}
