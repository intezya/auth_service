package persistence

import (
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
)

type Provider struct {
	AccountRepository repository.AccountRepository
}

func NewProvider(client *ent.Client) *Provider {
	return &Provider{AccountRepository: NewAccountRepositoryWithTracing(NewAccountRepository(client))}
}
