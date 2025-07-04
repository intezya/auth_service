// Code generated by tracing-gen. DO NOT EDIT.

package persistence

import (
	"context"
	domain "github.com/intezya/auth_service/internal/domain/account"
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
	tracer "github.com/intezya/auth_service/pkg/tracer"
)

type accountRepositoryWithTracing struct {
	wrapped repository.AccountRepository
}

func NewAccountRepositoryWithTracing(client *ent.Client) repository.AccountRepository {
	wrapped := NewAccountRepository(client)
	return &accountRepositoryWithTracing{
		wrapped: wrapped,
	}
}

func (t *accountRepositoryWithTracing) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	ctx, span := tracer.StartSpan(ctx, "AccountRepository.Create")
	defer span.End()

	return t.wrapped.Create(ctx, account)
}

func (t *accountRepositoryWithTracing) FindByID(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	ctx, span := tracer.StartSpan(ctx, "AccountRepository.FindByID")
	defer span.End()

	return t.wrapped.FindByID(ctx, id)
}

func (t *accountRepositoryWithTracing) FindByLowerUsername(ctx context.Context, username domain.Username) (*domain.Account, error) {
	ctx, span := tracer.StartSpan(ctx, "AccountRepository.FindByLowerUsername")
	defer span.End()

	return t.wrapped.FindByLowerUsername(ctx, username)
}

func (t *accountRepositoryWithTracing) ExistsByLowerUsername(ctx context.Context, username domain.Username) bool {
	ctx, span := tracer.StartSpan(ctx, "AccountRepository.ExistsByLowerUsername")
	defer span.End()

	return t.wrapped.ExistsByLowerUsername(ctx, username)
}

func (t *accountRepositoryWithTracing) Update(ctx context.Context, account *domain.Account) error {
	ctx, span := tracer.StartSpan(ctx, "AccountRepository.Update")
	defer span.End()

	return t.wrapped.Update(ctx, account)
}
