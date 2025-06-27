package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/intezya/auth_service/internal/adapters/http"
	authpb "github.com/intezya/auth_service/protos/go/auth"
	"net"
	"sync"

	"github.com/intezya/pkglib/logger"
	"google.golang.org/grpc"
)

var ErrShutdownTimeout = errors.New("shutdown timed out")

type App struct {
	server   *grpc.Server
	port     int
	listener net.Listener
	mu       sync.Mutex
	running  bool
}

func NewGRPCApp(provider *Provider, config Config) *App {
	server := grpc.NewServer()

	authpb.RegisterAuthServiceServer(server, provider.AuthController)

	http.SetupMetricsServer(config.MetricsPort)

	return &App{
		server: server,
		port:   config.GRPCServerPort,
	}
}

func (a *App) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return errors.New("server is already running")
	}
	a.running = true
	a.mu.Unlock()

	logger.Log.Infof("Starting gRPC server on port %d", a.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		return fmt.Errorf("failed to listen on port %d: %w", a.port, err)
	}

	a.mu.Lock()
	a.listener = lis
	a.mu.Unlock()

	logger.Log.Infof("gRPC server listening on %s", lis.Addr().String())

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		if err := a.server.Serve(lis); err != nil {
			select {
			case errCh <- fmt.Errorf("gRPC server serve error: %w", err):
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Infof("gRPC server context cancelled, stopping...")
		return ctx.Err()
	case err := <-errCh:
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		if err != nil {
			return err
		}
		return nil
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	a.mu.Lock()
	if !a.running {
		a.mu.Unlock()
		return nil
	}
	a.mu.Unlock()

	logger.Log.Infof("Shutting down gRPC server on port %d...", a.port)

	done := make(chan struct{})

	go func() {
		defer close(done)
		a.server.GracefulStop()

		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
	}()

	select {
	case <-done:
		logger.Log.Infof("gRPC server on port %d gracefully stopped", a.port)
		return nil
	case <-ctx.Done():
		logger.Log.Warnf("gRPC server on port %d shutdown timeout, forcing stop", a.port)

		a.server.Stop()

		a.mu.Lock()
		a.running = false
		a.mu.Unlock()

		select {
		case <-done:
		case <-ctx.Done():
		}

		return ErrShutdownTimeout
	}
}

func (a *App) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.running
}

func (a *App) Port() int {
	return a.port
}

func (a *App) Server() *grpc.Server {
	return a.server
}
