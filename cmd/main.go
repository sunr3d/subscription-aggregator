package main

import (
	"log"

	"github.com/sunr3d/subscription-aggregator/internal/config"
	"github.com/sunr3d/subscription-aggregator/internal/entrypoint"
	"github.com/sunr3d/subscription-aggregator/internal/logger"
)

func main() {
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatalf("Ошибка при получении конфигурации: %s\n", err.Error())
	}

	zapLogger := logger.New(cfg.LogLevel)

	if err = entrypoint.Run(cfg, zapLogger); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %s\n", err.Error())
	}
}
