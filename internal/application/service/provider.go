package service

import (
	"github.com/intezya/auth_service/internal/infrastructure/persistence"
)

type Provider struct {
	AuthService AuthService
}

func NewProvider(
	provider *persistence.Provider,
	credentialsHelper CredentialsHelper,
	tokenHelper TokenHelper,
) *Provider {
	return &Provider{
		AuthService: NewAuthServiceWithTracing(
			provider.AccountRepository,
			credentialsHelper,
			tokenHelper,
		),
	}
}
