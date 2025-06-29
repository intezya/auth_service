package grpc

import (
	"github.com/intezya/auth_service/internal/application/service"
	authpb "github.com/intezya/auth_service/protos/go/auth"
)

type Provider struct {
	AuthController authpb.AuthServiceServer
}

func NewProvider(provider *service.Provider) *Provider {
	return &Provider{
		AuthController: NewAuthControllerWithTracing(provider.AuthService),
	}
}
