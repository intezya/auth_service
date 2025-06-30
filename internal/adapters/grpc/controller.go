package grpc

import (
	"context"
	"github.com/intezya/auth_service/internal/application/usecase"
	authpb "github.com/intezya/auth_service/protos/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authController struct {
	authpb.UnimplementedAuthServiceServer

	authService usecase.AuthUseCase
}

func NewAuthController(authService usecase.AuthUseCase) authpb.AuthServiceServer {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(
	ctx context.Context,
	request *authpb.AuthenticationRequest,
) (*authpb.Empty, error) {
	if request.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if request.GetHardwareId() == "" {
		return nil, status.Error(codes.InvalidArgument, "hardware_id is required")
	}

	err := c.authService.Register(
		ctx, &usecase.RegisterCommand{
			Username:   request.Username,
			Password:   request.Password,
			HardwareID: request.HardwareId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.Empty{}, nil
}

func (c *authController) Login(
	ctx context.Context,
	request *authpb.AuthenticationRequest,
) (*authpb.TokenResponse, error) {
	if request.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if request.GetHardwareId() == "" {
		return nil, status.Error(codes.InvalidArgument, "hardware_id is required")
	}

	result, err := c.authService.Login(
		ctx, &usecase.LoginCommand{
			Username:   request.Username,
			Password:   request.Password,
			HardwareID: request.HardwareId,
		},
	)
	if err != nil {
		return nil, err
	}

	var bannedUntil int64 = 0
	if result.BannedUntil != nil {
		bannedUntil = result.BannedUntil.Unix()
	}

	return &authpb.TokenResponse{
		Token:             result.Token,
		AccessLevel:       int64(result.AccessLevel),
		BannedUntilInUnix: bannedUntil,
		IsBanned:          bannedUntil == 0,
	}, nil
}

func (c *authController) VerifyToken(
	ctx context.Context,
	request *authpb.VerifyTokenRequest,
) (*authpb.VerifyTokenResponse, error) {
	if request.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	result, err := c.authService.VerifyToken(
		ctx, &usecase.VerifyTokenCommand{
			Token: request.Token,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.VerifyTokenResponse{
		Subject:     int64(result.Subject),
		AccessLevel: int64(result.AccessLevel),
	}, nil
}

func (c *authController) BanAccount(ctx context.Context, request *authpb.BanAccountRequest) (*authpb.Empty, error) {
	if request.GetSubject() == 0 {
		return nil, status.Error(codes.InvalidArgument, "subject is required")
	}

	reason := &request.Reason
	if request.Reason == "" {
		reason = nil
	}

	err := c.authService.BanAccount(
		ctx,
		&usecase.BanAccountCommand{
			AccountID:    int(request.Subject),
			BanUntilUnix: request.BanUntilUnix,
			BanReason:    reason,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.Empty{}, nil
}
