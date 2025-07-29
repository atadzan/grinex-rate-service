package server

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/atadzan/grinex-rate-service/proto/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/atadzan/grinex-rate-service/internal/config"
	"github.com/atadzan/grinex-rate-service/internal/database"
	"github.com/atadzan/grinex-rate-service/internal/service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type RateServiceServer struct {
	pb.UnimplementedRateServiceServer
	db        *database.Database
	grinexSvc *service.GrinexService
	config    *config.Config
	logger    *zap.Logger
}

func NewRateServiceServer(cfg *config.Config, logger *zap.Logger) (*RateServiceServer, error) {
	db, err := database.NewDatabase(cfg.Database.GetDSN(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := database.RunMigrations(cfg.Database.GetDSN()); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	grinexConfig := &service.GrinexConfig{
		BaseURL:   cfg.Grinex.BaseURL,
		Timeout:   cfg.Grinex.Timeout,
		UserAgent: cfg.Grinex.UserAgent,
	}
	grinexSvc := service.NewGrinexService(grinexConfig, logger)

	return &RateServiceServer{
		db:        db,
		grinexSvc: grinexSvc,
		config:    cfg,
		logger:    logger,
	}, nil
}

func (s *RateServiceServer) GetRates(ctx context.Context, req *pb.GetRatesReq) (*pb.GetRatesResp, error) {
	ctx, span := otel.Tracer("grinex-rate-service").Start(ctx, "GetRates")
	defer span.End()

	s.logger.Info("GetRates called")

	rate, err := s.grinexSvc.GetUSDTRate(ctx)
	if err != nil {
		s.logger.Error("Failed to get rate from Grinex", zap.Error(err))
		return nil, fmt.Errorf("failed to get rate from Grinex: %w", err)
	}

	dbRecord := &database.RateRecord{
		TradingPair: rate.TradingPair,
		AskPrice:    rate.AskPrice,
		BidPrice:    rate.BidPrice,
		Timestamp:   rate.Timestamp,
		CreatedAt:   time.Now(),
	}

	if err := s.db.SaveRate(dbRecord); err != nil {
		s.logger.Error("Failed to save rate to database", zap.Error(err))
		return nil, fmt.Errorf("failed to save rate to database: %w", err)
	}

	// Convert to protobuf response
	return &pb.GetRatesResp{
		TradingPair: rate.TradingPair,
		AskPrice:    rate.AskPrice,
		BidPrice:    rate.BidPrice,
		Timestamp:   timestamppb.New(rate.Timestamp),
	}, nil
}

func (s *RateServiceServer) Healthcheck(ctx context.Context, req *pb.HealthcheckReq) (*pb.HealthcheckResp, error) {
	ctx, span := otel.Tracer("grinex-rate-service").Start(ctx, "Healthcheck")
	defer span.End()

	s.logger.Info("Healthcheck called")

	status := "healthy"
	message := "Service is healthy"

	// Check database health
	if err := s.db.HealthCheck(); err != nil {
		status = "unhealthy"
		message = fmt.Sprintf("Database health check failed: %v", err)
		s.logger.Error("Database health check failed", zap.Error(err))
		return &pb.HealthcheckResp{
			Status:  status,
			Message: message,
		}, fmt.Errorf("database health check failed: %w", err)
	}

	// Check Grinex API health
	if err := s.grinexSvc.HealthCheck(ctx); err != nil {
		status = "degraded"
		message = fmt.Sprintf("Grinex API health check failed: %v", err)
		s.logger.Warn("Grinex API health check failed", zap.Error(err))
	}

	return &pb.HealthcheckResp{
		Status:  status,
		Message: message,
	}, nil
}

func (s *RateServiceServer) Close() error {
	return s.db.Close()
}

func StartServer(ctx context.Context, cfg *config.Config, logger *zap.Logger) error {
	server, err := NewRateServiceServer(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	defer server.Close()

	port := ":" + cfg.Server.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRateServiceServer(s, server)

	reflection.Register(s)

	logger.Info("gRPC server listening", zap.String("port", port))

	// Start server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("Server error", zap.Error(err))
		}
	}()

	// Wait for context cancellation (graceful shutdown)
	<-ctx.Done()

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.GracefulStop()

	// Wait for shutdown to complete
	<-shutdownCtx.Done()

	logger.Info("Server stopped gracefully")
	return nil
}

// SetupMetrics sets up Prometheus metrics
func SetupMetrics() (*metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	return provider, nil
}
