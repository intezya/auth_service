package main

import (
	"context"
	"fmt"
	"github.com/intezya/auth_service/internal/adapters/config"
	"github.com/intezya/auth_service/internal/adapters/grpc"
	"github.com/intezya/auth_service/internal/application/usecase"
	"github.com/intezya/auth_service/internal/domain/service"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
	"github.com/intezya/auth_service/internal/infrastructure/persistence"
	"github.com/intezya/auth_service/internal/pkg/crypto"
	"github.com/intezya/auth_service/internal/pkg/jwt"
	domainvalidator "github.com/intezya/auth_service/internal/pkg/validator"
	"github.com/intezya/auth_service/pkg/tracer"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/intezya/pkglib/errorz"
	"github.com/intezya/pkglib/logger"
)

const gracefulShutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		logger.Log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	config := config.LoadConfig()

	_, err := logger.New(
		logger.WithCaller(config.Logger.CallerEnabled),
		logger.WithDebug(config.Logger.Debug),
		logger.WithEnvironment(config.Logger.Environment),
		logger.WithTimeZone(config.Logger.TimeZone),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := tracer.Init(config.Tracer, logger.Log); err != nil {
		logger.Log.Warnf("Failed to initialize tracer: %v", err)
	}

	errorz.SetValidator(validator.New())
	validators := domainvalidator.NewProvider()
	tokenManager := jwt.NewTokenManager(config.JWT)
	passwordEncoder := crypto.NewPasswordEncoder(config.Crypto)
	entClient := persistence.SetupEnt(config.Ent, logger.Log)

	repositories := persistence.NewProvider(entClient)
	hardwareIDManager := service.NewHardwareIDManager(repositories.AccountRepository, passwordEncoder)
	services := usecase.NewProvider(repositories, validators, passwordEncoder, tokenManager, hardwareIDManager)
	controllers := grpc.NewProvider(services)
	grpcApp := grpc.NewGRPCApp(controllers, config.Server)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcApp.Start(ctx); err != nil {
			logger.Log.Errorf("gRPC server error: %v", err)
			stop()
		}
	}()

	logger.Log.Info("Application started successfully")

	<-ctx.Done()
	logger.Log.Info("Shutdown signal received")

	return gracefulShutdown(grpcApp, entClient, &wg)
}

func gracefulShutdown(grpcApp *grpc.App, entClient *ent.Client, wg *sync.WaitGroup) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := grpcApp.Shutdown(shutdownCtx); err != nil {
		logger.Log.Errorf("Shutdown error: %v", err)
	}

	wg.Wait()

	if err := tracer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Warnf("Tracer shutdown error: %v", err)
	}

	if err := entClient.Close(); err != nil {
		logger.Log.Warnf("Ent client close error: %v", err)
	}

	logger.Log.Info("Application shutdown completed")
	return nil
}
