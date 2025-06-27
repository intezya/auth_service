package main

import (
	"context"
	"fmt"
	"github.com/intezya/auth_service/internal/adapters/config"
	"github.com/intezya/auth_service/internal/adapters/grpc"
	"github.com/intezya/auth_service/internal/application/service"
	"github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
	"github.com/intezya/auth_service/internal/infrastructure/persistence"
	"github.com/intezya/auth_service/internal/pkg/crypto"
	"github.com/intezya/auth_service/internal/pkg/jwt"
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

	tracer.Init(config.Tracer, logger.Log)
	errorz.SetValidator(validator.New())

	tokenHelper := jwt.NewTokenHelper(config.JWT)
	hashHelper := crypto.NewHashHelper(config.Crypto)
	entClient := persistence.SetupEnt(config.Ent, logger.Log)

	repositories := persistence.NewProvider(entClient)
	services := service.NewProvider(repositories, hashHelper, tokenHelper)
	controllers := grpc.NewProvider(services)

	grpcApp := grpc.NewGRPCApp(controllers, config.Server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	var wg sync.WaitGroup

	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcApp.Start(ctx); err != nil {
			select {
			case errCh <- fmt.Errorf("gRPC server error: %w", err):
			default:
			}
		}
	}()

	logger.Log.Info("Application started successfully")

	select {
	case sig := <-sigCh:
		logger.Log.Infof("Received shutdown signal: %v", sig)
	case err := <-errCh:
		logger.Log.Errorf("Service error: %v", err)
		cancel()
	}

	// Graceful shutdown
	logger.Log.Info("Starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		gracefulShutdownTimeout,
	)
	defer shutdownCancel()

	cancel()

	shutdownDone := make(chan error, 1)
	go func() {
		if err := grpcApp.Shutdown(shutdownCtx); err != nil {
			shutdownDone <- err
			return
		}
		shutdownDone <- nil
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			logger.Log.Errorf("Shutdown error: %v", err)
		}
	case <-shutdownCtx.Done():
		logger.Log.Warn("Shutdown timeout exceeded")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Info("All services stopped")
	case <-time.After(time.Second):
		logger.Log.Warn("Some services may still be running")
	}

	logger.Log.Info("Application shutdown completed")
	return nil
}
