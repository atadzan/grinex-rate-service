package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/atadzan/grinex-rate-service/internal/config"
	"github.com/atadzan/grinex-rate-service/server"
)

func main() {
	cfg := config.Load()

	logger, err := initLogger(cfg.Logging.Level)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	_, err = server.SetupMetrics()
	if err != nil {
		logger.Error("Failed to setup metrics", zap.Error(err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handling graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
	}()

	logger.Info("Starting gRPC Rate Service",
		zap.String("port", cfg.Server.Port),
		zap.String("database_host", cfg.Database.Host),
		zap.String("grinex_base_url", cfg.Grinex.BaseURL),
	)

	if err := server.StartServer(ctx, cfg, logger); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initLogger(level string) (*zap.Logger, error) {
	var logLevel zap.AtomicLevel
	switch level {
	case "debug":
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config := zap.NewProductionConfig()
	config.Level = logLevel

	return config.Build()
}
