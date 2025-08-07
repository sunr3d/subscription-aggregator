package entrypoint

import (
	"go.uber.org/zap"

	"github.com/sunr3d/subscription-aggregator/internal/config"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	// TODO: Сборка зависимостей, запуск сервиса, старт http сервера
	return nil
}