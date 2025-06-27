package grpc

import (
	"github.com/intezya/auth_service/internal/application/service"
)

type Provider struct {
	AuthController *AuthController
}

func NewProvider(provider *service.Provider) *Provider {
	return &Provider{
		AuthController: NewAuthController(provider.AuthService),
	}
}
