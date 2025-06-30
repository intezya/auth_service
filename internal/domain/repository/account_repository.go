package repository

import (
	"context"
	domain "github.com/intezya/auth_service/internal/domain/account"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	FindByID(ctx context.Context, id domain.AccountID) (*domain.Account, error)
	FindByLowerUsername(ctx context.Context, username domain.Username) (*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	ExistsByLowerUsername(ctx context.Context, username domain.Username) bool
}
