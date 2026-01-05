package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/harpy-wings/mastering-go/s0e0/internal/usermanager"
	"github.com/harpy-wings/mastering-go/s0e0/internal/walletmanager"
	"github.com/harpy-wings/mastering-go/s0e0/pkg/constants"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type oracleApp struct {
	logger *logrus.Logger
	rdb    redis.UniversalClient

	httpServer *http.Server

	// userManager   usermanager.IUserManager
	// walletManager walletmanager.IWalletManager

	onStop []func(ctx context.Context) error
}

var _ IApp = &oracleApp{}

func (a *oracleApp) Run(ctx context.Context) error {
	err := a.httpServer.ListenAndServe()

	<-ctx.Done()
	a.logger.Info("http server stopped")
	if err != nil && err != http.ErrServerClosed {
		a.logger.WithError(err).Error("failed to listen and serve http server")
		return err
	}

	return nil
}

func (a *oracleApp) GracefulStop(ctx context.Context) error {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.logger.WithError(err).Error("failed to shutdown http server")
	}

	for i := len(a.onStop) - 1; i >= 0; i-- {
		if err := a.onStop[i](ctx); err != nil {
			a.logger.WithError(err).Error("failed to stop onStop function")
		}
	}

	return nil
}

func NewOracleApp(ctx context.Context, v *viper.Viper) (IApp, error) {
	a := new(oracleApp)
	err := a.presets(ctx, v)
	if err != nil {
		return nil, err
	}
	err = a.init(ctx, v)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *oracleApp) presets(ctx context.Context, v *viper.Viper) error {
	// Init Logger
	{
		logger := logrus.StandardLogger()
		logger.SetLevel(logrus.DebugLevel)
		logger.ReportCaller = v.GetBool("logger.report_caller")
		{
			if lgLevel := v.GetString("logger.log_level"); lgLevel != "" {
				lv, err := logrus.ParseLevel(lgLevel)
				if err != nil {
					return err
				}
				logger.SetLevel(lv)
			}
		}
		if formatter := v.GetString("logger.formatter"); formatter != "" {
			switch formatter {
			case "json":
				logger.SetFormatter(&logrus.JSONFormatter{})
			case "text":
				logger.SetFormatter(&logrus.TextFormatter{})
			default:
				return fmt.Errorf("invalid logger formatter: %s", formatter)
			}
		}
		if output := v.GetString("logger.output"); output != "" {
			switch output {
			case "stdout":
				logger.SetOutput(os.Stdout)
			default:
				output = filepath.Join("/var/log/", filepath.Clean(output))
				if !strings.HasPrefix(output, "/var/log/") {
					return fmt.Errorf("invalid path")
				}
				file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
				if err != nil {
					return err
				}
				logger.SetOutput(file)
			}
		}
		a.logger = logger
	}
	// Init Sentry
	// Init Instrumentation
	// Init Tracing

	// Init Redis
	{
		a.logger.Info("initializing redis")

		// Check if we have multiple addresses (cluster mode) or single address (single instance)
		addrs := v.GetStringSlice("redis.url")
		if len(addrs) > 1 {
			// Use cluster client for multiple addresses

			for i := range addrs {
				addrs[i] = strings.TrimPrefix(addrs[i], "redis://")
				addrs[i] = strings.TrimPrefix(addrs[i], "rediss://")
			}
			a.logger.Info("using redis cluster client")
			rdb := redis.NewClusterClient(&redis.ClusterOptions{
				ClientName: strings.Join([]string{v.GetString("app.environment"), v.GetString("app.name"), constants.VERSION}, "-"),
				Addrs:      addrs,
				Username:   v.GetString("redis.username"),
				Password:   v.GetString("redis.password"),
			})
			a.rdb = rdb
		} else {
			// Use single client for single address
			a.logger.Info("using redis single client")
			addr := "localhost:6379" // default
			if len(addrs) > 0 {
				addr = addrs[0]
			}

			// Parse the address to handle redis:// protocol
			addr = strings.TrimPrefix(addr, "redis://")

			rdb := redis.NewClient(&redis.Options{
				Addr:     addr,
				Username: v.GetString("redis.username"),
				Password: v.GetString("redis.password"),
				DB:       v.GetInt("redis.db"),
			})
			a.rdb = rdb
		}

		err := a.rdb.Ping(ctx).Err()
		if err != nil {
			return err
		}
		a.onStop = append(a.onStop, func(ctx context.Context) error {
			a.logger.Info("1.closing redis")
			return a.rdb.Close()
		})
	}

	return nil
}

func (a *oracleApp) init(ctx context.Context, v *viper.Viper) error {
	// Init HTTP Server
	{
		a.logger.Info("initializing http server")
		a.httpServer = &http.Server{
			Addr: fmt.Sprintf(":%d", v.GetInt("http.port")),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello, Go Masters!"))
			}),
		}
	}
	// Init User Manager

	var UM usermanager.IUserManager

	{
		WM, err := walletmanager.New(ctx,
			walletmanager.WithLogger(a.logger.WithField("module", "walletmanager")),
			walletmanager.WithRedis(a.rdb),
			walletmanager.WithUserManager(UM),
		)
		if err != nil {
			return err
		}
		a.onStop = append(a.onStop, func(ctx context.Context) error {
			return WM.GracefulStop(ctx)
		})
	}

	return nil
}
