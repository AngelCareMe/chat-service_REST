package app

import (
	"context"
	"net/http"
	"time"

	postgres "chat-service/internal/adapter"
	"chat-service/internal/handler"

	"github.com/sirupsen/logrus"
)

type App struct {
	httpServer *http.Server
	dbAdapter  *postgres.PostgresAdapter
	handler    *handler.Handler
	logger     *logrus.Logger
}

func NewApp(
	httpServer *http.Server,
	dbAdapter *postgres.PostgresAdapter,
	handler *handler.Handler,
	logger *logrus.Logger,
) *App {
	return &App{
		httpServer: httpServer,
		dbAdapter:  dbAdapter,
		handler:    handler,
		logger:     logger,
	}
}

func (a *App) Start() error {
	a.logger.Info("starting application server")

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.WithError(err).Fatal("failed to start server")
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("shutting down application gracefully")

	// Закрываем HTTP сервер с таймаутом
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.WithError(err).Error("failed to shutdown HTTP server gracefully")
		// Принудительное закрытие
		a.httpServer.Close()
	}

	// Закрываем соединение с БД
	if a.dbAdapter != nil {
		a.dbAdapter.Close()
	}

	// Закрываем handler ресурсы
	if a.handler != nil {
		a.handler.Close()
	}

	a.logger.Info("application shutdown completed")
	return nil
}
