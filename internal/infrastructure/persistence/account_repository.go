package persistence

import (
	"context"
	"github.com/intezya/auth_service/internal/adapters/mapper"
	domain "github.com/intezya/auth_service/internal/domain/account"
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
	entAccount "github.com/intezya/auth_service/internal/infrastructure/ent/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type accountRepository struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) repository.AccountRepository {
	return &accountRepository{client: client}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	created, err := r.client.Account.
		Create().
		SetUsername(account.Username()).
		SetPassword(account.Password()).
		SetNillableHardwareID(account.HardwareID()).
		Save(ctx)
	if err != nil {
		return nil, r.handleConstraintError(err)
	}

	return mapper.EntAccountToDomain(created), nil
}

func (r *accountRepository) FindByID(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	found, err := r.client.Account.Get(ctx, int(id))
	if err != nil {
		return nil, r.handleNotFoundError(err)
	}

	return mapper.EntAccountToDomain(found), nil
}

func (r *accountRepository) FindByLowerUsername(ctx context.Context, username domain.Username) (
	*domain.Account,
	error,
) {
	found, err := r.client.Account.
		Query().
		Where(entAccount.UsernameEqualFold(string(username))).
		First(ctx)
	if err != nil {
		return nil, r.handleNotFoundError(err)
	}

	return mapper.EntAccountToDomain(found), nil
}

func (r *accountRepository) ExistsByLowerUsername(ctx context.Context, username domain.Username) bool {
	exists, err := r.client.Account.
		Query().
		Where(entAccount.UsernameEqualFold(string(username))).
		Exist(ctx)

	return err != nil && exists
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) error {
	// update can return not found, but in code repository.Update called only if account already found
	err := r.client.Account.UpdateOne(
		&ent.Account{
			ID:          account.ID(),
			Username:    account.Username(),
			Password:    account.Password(),
			HardwareID:  account.HardwareID(),
			AccessLevel: domain.AccessLevel(account.AccessLevel()),
			CreatedAt:   account.CreatedAt(),
			BannedUntil: account.BannedUntil(),
			BanReason:   account.BanReason(),
		},
	).Exec(ctx)
	if err != nil {
		return r.handleConstraintError(err)
	}

	return nil
}

func (r *accountRepository) handleConstraintError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case strings.Contains(err.Error(), "username"):
		return status.Errorf(codes.AlreadyExists, "user already exists")
	case strings.Contains(err.Error(), "hardware_id"):
		return status.Errorf(codes.AlreadyExists, "hardware_id conflict")
	}

	return status.Errorf(codes.Internal, "unexpected internal error: %v", err)
}

func (r *accountRepository) handleNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return status.Errorf(codes.NotFound, "account not found")
	}

	return status.Errorf(codes.Internal, "unexpected internal error: %v", err)
}
