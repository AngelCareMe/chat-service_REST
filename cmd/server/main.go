// @title Chat Service API
// @version 1.0
// @description REST API для чат-приложения
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer YOUR_TOKEN' to authenticate"
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgres "chat-service/internal/adapter"
	"chat-service/internal/app"
	"chat-service/internal/handler"
	"chat-service/internal/service"
	"chat-service/internal/usecase/message"
	"chat-service/internal/usecase/session"
	"chat-service/internal/usecase/user"
	"chat-service/pkg/config"
	"chat-service/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// @title Chat Service API
// @version 1.0
// @description REST API для чат-приложения
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer YOUR_TOKEN' to authenticate"
func main() {
	// Initialize logger first
	appLogger := logger.NewLogger()
	appLogger.Info("Starting chat service application")

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		appLogger.WithError(err).Fatal("failed to load configuration")
	}

	// Set logger level from config
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		appLogger.WithError(err).Warnf("invalid log level %s, using info", cfg.Logger.Level)
		level = logrus.InfoLevel
	}
	appLogger.SetLevel(level)

	// Print configuration
	cfg.Print()

	// Initialize database connection pool
	dbPool, err := initDatabase(cfg, appLogger)
	if err != nil {
		appLogger.WithError(err).Fatal("failed to initialize database")
	}
	defer func() {
		appLogger.Info("closing database connection pool")
		dbPool.Close()
	}()

	// Initialize adapters
	dbAdapter := postgres.NewPostgresAdapter(dbPool, appLogger)

	// Initialize services
	hashService := service.NewHashService(appLogger)
	jwtService := service.NewJWTService(cfg.JWT.SecretKey, appLogger)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbAdapter)
	messageRepo := postgres.NewMessageRepository(dbAdapter)
	sessionRepo := postgres.NewSessionRepository(dbAdapter)

	// Initialize usecases
	userUsecase := user.NewUserUsecase(userRepo, sessionRepo, hashService, jwtService, appLogger)
	messageUsecase := message.NewMessageUsecase(messageRepo, userRepo, appLogger)
	sessionUsecase := session.NewSessionUsecase(sessionRepo, jwtService, appLogger)

	// Initialize handler
	appHandler := handler.NewHandler(userUsecase, messageUsecase, sessionUsecase, appLogger)

	// Initialize HTTP server
	httpServer := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      appHandler.GetRouter(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Create application instance
	application := app.NewApp(httpServer, dbAdapter, appHandler, appLogger)

	// Start server in a goroutine
	appLogger.WithField("address", cfg.GetServerAddress()).Info("starting HTTP server")
	go func() {
		if err := application.Start(); err != nil && err != http.ErrServerClosed {
			appLogger.WithError(err).Fatal("failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform graceful shutdown
	if err := application.Stop(ctx); err != nil {
		appLogger.WithError(err).Fatal("failed to shutdown application gracefully")
	}

	appLogger.Info("server exited gracefully")
}

// initDatabase initializes the database connection pool
func initDatabase(cfg *config.Config, logger *logrus.Logger) (*pgxpool.Pool, error) {
	logger.Info("initializing database connection pool")

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = cfg.Database.MaxConnections
	poolConfig.MinConns = cfg.Database.MinConnections
	poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection pool initialized successfully")
	return pool, nil
}
