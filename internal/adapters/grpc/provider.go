package grpc

import (
	"github.com/intezya/auth_service/internal/application/usecase"
	authpb "github.com/intezya/auth_service/protos/go/auth"
)

type Provider struct {
	AuthController authpb.AuthServiceServer
}

func NewProvider(provider *usecase.Provider) *Provider {
	return &Provider{
		AuthController: NewAuthControllerWithTracing(provider.AuthUseCase),
	}
}
