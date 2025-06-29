package persistence

import (
	"context"
	"github.com/intezya/auth_service/internal/adapters/mapper"
	"github.com/intezya/auth_service/internal/domain/dto"
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
	entAccount "github.com/intezya/auth_service/internal/infrastructure/ent/account"
	"github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type accountRepository struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) repository.AccountRepository {
	return &accountRepository{client: client}
}

func (r *accountRepository) Create(
	ctx context.Context,
	username string,
	password string,
	hardwareId string,
) (*dto.AccountDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "accountRepository.Create")
	defer span.End()

	account, err := r.client.Account.
		Create().
		SetUsername(username).
		SetPassword(password).
		SetHardwareID(hardwareId).
		Save(ctx)
	if err != nil {
		return nil, r.handleConstraintError(err)
	}

	return mapper.AccountToDto(account), nil
}

func (r *accountRepository) FindByID(ctx context.Context, id int) (*dto.AccountDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "accountRepository.FindByID")
	defer span.End()

	account, err := r.client.Account.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "unexpected internal error: %v", err)
	}

	return mapper.AccountToDto(account), nil
}

func (r *accountRepository) FindByLowerUsername(ctx context.Context, username string) (*dto.AccountDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "accountRepository.FindByLowerUsername")
	defer span.End()

	account, err := r.client.Account.
		Query().
		Where(entAccount.UsernameEqualFold(username)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "unexpected internal error: %v", err)
	}

	return mapper.AccountToDto(account), nil
}

func (r *accountRepository) UpdateHardwareIDByID(ctx context.Context, id int, hardwareId string) error {
	ctx, span := tracer.StartSpan(ctx, "accountRepository.UpdateHardwareIDByID")
	defer span.End()

	_, err := r.client.Account.
		UpdateOneID(id).
		SetHardwareID(hardwareId).
		Save(ctx)

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

func (r *accountRepository) UpdateBannedUntilBannedReasonByID(
	ctx context.Context,
	id int,
	bannedUntil *time.Time,
	banReason *string,
) error {
	_, err := r.client.Account.
		UpdateOneID(id).
		SetNillableBannedUntil(bannedUntil).
		SetNillableBanReason(banReason).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return status.Errorf(codes.NotFound, "account not found")
		}

		return status.Errorf(codes.Internal, "unexpected internal error: %v", err)
	}

	return nil
}
