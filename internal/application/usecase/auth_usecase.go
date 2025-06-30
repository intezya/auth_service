package usecase

import (
	"context"
	entity "github.com/intezya/auth_service/internal/domain/account"
	"github.com/intezya/auth_service/internal/domain/dto"
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/domain/service"
	"github.com/intezya/auth_service/pkg/clock"
	"time"
)

type AuthUseCase interface {
	Register(ctx context.Context, cmd *RegisterCommand) error
	Login(ctx context.Context, cmd *LoginCommand) (*LoginResult, error)
	VerifyToken(ctx context.Context, cmd *VerifyTokenCommand) (*dto.TokenData, error)
	BanAccount(ctx context.Context, cmd *BanAccountCommand) error
}

type RegisterCommand struct {
	Username   string
	Password   string
	HardwareID string
}

type LoginCommand struct {
	Username   string
	Password   string
	HardwareID string
}

type VerifyTokenCommand struct {
	Token string
}

type BanAccountCommand struct {
	AccountID    int
	BanUntilUnix int64
	BanReason    *string
}

type LoginResult struct {
	Token       string
	AccessLevel int
	BannedUntil *time.Time // nil = not banned
}

type authUseCase struct {
	accountRepository repository.AccountRepository

	tokenManager service.TokenManager

	passwordValidator service.Validator[string]
	usernameValidator service.Validator[string]
	hardwareValidator service.Validator[string]
	passwordEncoder   service.PasswordEncoder
	hardwareIDManager service.HardwareIDManager

	clock clock.Clock
}

func NewAuthUseCase(
	accountRepository repository.AccountRepository,
	passwordEncoder service.PasswordEncoder,
	tokenManager service.TokenManager,
	hardwareIDManager service.HardwareIDManager,
	usernameValidator service.Validator[string],
	passwordValidator service.Validator[string],
	hardwareValidator service.Validator[string],
) AuthUseCase {
	return &authUseCase{
		accountRepository: accountRepository,
		passwordEncoder:   passwordEncoder,
		tokenManager:      tokenManager,
		usernameValidator: usernameValidator,
		passwordValidator: passwordValidator,
		hardwareValidator: hardwareValidator,
		hardwareIDManager: hardwareIDManager,
	}
}

func (uc *authUseCase) Register(ctx context.Context, cmd *RegisterCommand) error {
	err := uc.passwordValidator.Validate(cmd.Username)
	if err != nil {
		return err
	}

	err = uc.passwordValidator.Validate(cmd.Password)
	if err != nil {
		return err
	}

	err = uc.hardwareValidator.Validate(cmd.HardwareID)
	if err != nil {
		return err
	}

	if uc.accountRepository.ExistsByLowerUsername(ctx, entity.Username(cmd.Username)) {
		panic("TODO()")
		//return TODO()
	}

	encodedPassword := uc.passwordEncoder.EncodePassword(ctx, cmd.Password)
	encodedHardwareID := uc.passwordEncoder.EncodeHardwareID(ctx, cmd.HardwareID)

	newAccount := entity.NewAccount(
		entity.Username(cmd.Username),
		entity.HashedPassword(encodedPassword),
		entity.HardwareID(encodedHardwareID),
		uc.clock,
	)

	_, err = uc.accountRepository.Create(ctx, newAccount) // hardware id conflict

	return err
}

func (uc *authUseCase) Login(ctx context.Context, cmd *LoginCommand) (*LoginResult, error) {
	account, err := uc.accountRepository.FindByLowerUsername(ctx, entity.Username(cmd.Username))
	if err != nil {
		return nil, err
	}

	if !uc.passwordEncoder.VerifyPassword(ctx, cmd.Password, account.Password()) {
		panic("TODO()")
		//return nil, TODO()
	}

	err = uc.hardwareIDManager.ValidateAndSetHardwareID(ctx, account, cmd.HardwareID)
	if err != nil {
		return nil, err
	}

	if account.IsBanned(uc.clock) {
		panic("TODO()")
		//return nil, TODO()
	}

	token := uc.tokenManager.Generate(account.ID())

	return &LoginResult{
		Token:       token,
		AccessLevel: account.AccessLevel(),
		BannedUntil: account.BannedUntil(),
	}, nil

}

func (uc *authUseCase) VerifyToken(ctx context.Context, cmd *VerifyTokenCommand) (*dto.TokenData, error) {
	tokenData, err := uc.tokenManager.Parse(cmd.Token)
	if err != nil {
		return nil, err
	}

	account, err := uc.accountRepository.FindByID(ctx, entity.AccountID(tokenData.Subject))
	if err != nil {
		return nil, err
	}

	if account.IsBanned(uc.clock) {
		panic("TODO()")
		//return nil, TODO()
	}

	tokenData.AccessLevel = account.AccessLevel()

	return tokenData, nil
}

func (uc *authUseCase) BanAccount(ctx context.Context, cmd *BanAccountCommand) error {
	account, err := uc.accountRepository.FindByID(ctx, entity.AccountID(cmd.AccountID))
	if err != nil {
		return err
	}

	banUntilAsTime := uc.clock.Unix(cmd.BanUntilUnix, 0)

	if cmd.BanUntilUnix == 0 {
		account.Unban()
	} else {
		err = account.Ban(banUntilAsTime, cmd.BanReason)
		if err != nil {
			return err
		}
	}

	err = uc.accountRepository.Update(ctx, account)
	if err != nil {
		return err // only unexpected
	}

	return nil
}
