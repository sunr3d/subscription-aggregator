package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
	shutdownTimeout time.Duration
}

func New(port string, handler http.Handler, timeout time.Duration,logger *zap.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:              ":" + port,
			Handler:           handler,
			ReadTimeout:       timeout,
			WriteTimeout:      timeout,
			IdleTimeout:       timeout,
			ReadHeaderTimeout: timeout,
		},
		logger: logger,
		shutdownTimeout: timeout,
	}
}

func (s *Server) Start(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		s.logger.Info("Запуск HTTP сервера",
			zap.String("address", s.server.Addr),
		)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("ошибка HTTP сервера: %w", err)
			return
		}
		serverErr <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		
		s.logger.Info("Получен сигнал завершения")
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Ошибка при завершении работы сервера", zap.Error(err))
			return fmt.Errorf("ошибка при завершении работы сервера: %w", err)
		}
		s.logger.Info("Сервер остановлен корректно")
		return nil
	case err := <-serverErr:
		return err
	}
}
