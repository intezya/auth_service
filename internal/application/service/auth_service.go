package service

import (
	"context"
	"github.com/intezya/auth_service/internal/domain/dto"
	"github.com/intezya/auth_service/internal/domain/repository"
	"github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type authService struct {
	accountRepository repository.AccountRepository
	credentialsHelper CredentialsHelper
	tokenHelper       TokenHelper
}

func NewAuthService(
	accountRepository repository.AccountRepository,
	credentialsHelper CredentialsHelper,
	tokenHelper TokenHelper,
) AuthService {
	return &authService{
		accountRepository: accountRepository,
		credentialsHelper: credentialsHelper,
		tokenHelper:       tokenHelper,
	}
}

func (s *authService) Register(
	ctx context.Context,
	username string,
	password string,
	hardwareId string,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthService.Register")
	defer span.End()

	_, err := s.accountRepository.Create(
		ctx,
		username,
		s.encodePassword(ctx, password),
		s.encodeHardwareId(ctx, hardwareId),
	)

	return err
}

func (s *authService) Login(
	ctx context.Context,
	username string,
	password string,
	hardwareId string,
) (*dto.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthService.Login")
	defer span.End()

	account, err := s.accountRepository.FindByLowerUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !s.credentialsHelper.VerifyPassword(password, account.Password) {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	if account.HardwareID != nil {
		if !s.credentialsHelper.VerifyHardwareID(hardwareId, *account.HardwareID) {
			return nil, status.Error(codes.InvalidArgument, "invalid hardware_id")
		}
	} else {
		encoded := s.encodeHardwareId(ctx, hardwareId)
		account.HardwareID = &encoded
		err = s.accountRepository.UpdateHardwareIDByID(ctx, account.ID, encoded)
		if err != nil {
			return nil, err
		}
	}

	if s.checkAccountBanned(ctx, account) {
		return nil, status.Error(codes.PermissionDenied, "account banned")
	}

	return &dto.AuthenticationResult{
		AccessLevel: account.AccessLevel,
		Token:       s.tokenHelper.Generate(account.ID),
	}, nil
}

func (s *authService) VerifyToken(ctx context.Context, token string) (*dto.DataFromToken, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthService.VerifyToken")
	defer span.End()

	data, err := s.tokenHelper.Parse(token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	account, err := s.accountRepository.FindByID(ctx, data.Subject)
	if err != nil {
		return nil, err
	}

	if s.checkAccountBanned(ctx, account) {
		return nil, status.Error(codes.PermissionDenied, "account banned")
	}

	data.AccessLevel = account.AccessLevel

	return data, nil
}

func (s *authService) encodePassword(ctx context.Context, password string) string {
	ctx, span := tracer.StartSpan(ctx, "AuthService.encodePassword")
	defer span.End()

	return s.credentialsHelper.EncodePassword(password)
}

func (s *authService) encodeHardwareId(ctx context.Context, hardwareId string) string {
	ctx, span := tracer.StartSpan(ctx, "AuthService.encodeHardwareId")
	defer span.End()

	return s.credentialsHelper.EncodeHardwareID(hardwareId)
}

func (s *authService) BanAccount(ctx context.Context, subject int, banUntilUnix int, banReason string) error {
	ctx, span := tracer.StartSpan(ctx, "AuthService.BanAccount")
	defer span.End()

	reason := &banReason
	if banReason == "" {
		reason = nil
	}

	banUntilAsTime := time.Unix(int64(banUntilUnix), 0)

	return s.accountRepository.UpdateBannedUntilBannedReasonByID(ctx, subject, &banUntilAsTime, reason)
}

func (s *authService) checkAccountBanned(ctx context.Context, account *dto.AccountDTO) bool {
	if account.BannedUntil == nil {
		return false
	}

	if account.BannedUntil.Unix() <= time.Now().Unix() {
		_ = s.accountRepository.UpdateBannedUntilBannedReasonByID(
			ctx,
			account.ID,
			nil,
			nil,
		) // cannot be error if called with existing user
		return false
	}

	return true
}
