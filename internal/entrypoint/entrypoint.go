package entrypoint

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/sunr3d/subscription-aggregator/internal/api"
	"github.com/sunr3d/subscription-aggregator/internal/config"
	"github.com/sunr3d/subscription-aggregator/internal/infra/postgres"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/internal/middleware"
	"github.com/sunr3d/subscription-aggregator/internal/server"
	"github.com/sunr3d/subscription-aggregator/internal/services/subscription_service"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инфра слой
	db, err := postgres.New(cfg.Postgres, logger)
	if err != nil {
		return fmt.Errorf("postgres.New(): %w", err)
	}
	defer func(db infra.Database) {
		if c, ok := db.(interface{ Close() }); ok {
			c.Close()
		}
	}(db)

	// Сервисный слой
	svc := subscription_service.New(db)

	// API
	controller := api.New(svc, logger)
	mux := http.NewServeMux()
	controller.RegisterHandlers(mux)

	// Middleware
	handler := middleware.Recovery(logger)(
		middleware.ReqLogger(logger)(
			middleware.JSONValidator(logger)(mux),
		),
	)

	// HTTP сервер
	srv := server.New(cfg.HTTPPort, handler, cfg.HTTPTimeout, logger)

	return srv.Start(appCtx)
}
