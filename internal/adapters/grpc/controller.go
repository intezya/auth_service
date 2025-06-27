package grpc

import (
	"context"
	"github.com/intezya/auth_service/internal/application/service"
	"github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
	authpb "github.com/intezya/auth_service/protos/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	authpb.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(
	ctx context.Context,
	request *authpb.AuthenticationRequest,
) (*authpb.Empty, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthController.Register")
	defer span.End()

	if request.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if request.GetHardwareId() == "" {
		return nil, status.Error(codes.InvalidArgument, "hardware_id is required")
	}

	err := c.authService.Register(ctx, request.Username, request.Password, request.HardwareId)
	if err != nil {
		return nil, err
	}

	return &authpb.Empty{}, nil
}

func (c *AuthController) Login(
	ctx context.Context,
	request *authpb.AuthenticationRequest,
) (*authpb.TokenResponse, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthController.Login")
	defer span.End()

	if request.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if request.GetHardwareId() == "" {
		return nil, status.Error(codes.InvalidArgument, "hardware_id is required")
	}

	result, err := c.authService.Login(ctx, request.Username, request.Password, request.HardwareId)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResponse{
		Token:       result.Token,
		AccessLevel: int64(result.AccessLevel),
	}, nil
}

func (c *AuthController) VerifyToken(
	ctx context.Context,
	request *authpb.VerifyTokenRequest,
) (*authpb.VerifyTokenResponse, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthController.VerifyToken")
	defer span.End()

	if request.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	result, err := c.authService.VerifyToken(ctx, request.GetToken())
	if err != nil {
		return nil, err
	}

	return &authpb.VerifyTokenResponse{
		Subject:     int64(result.Subject),
		AccessLevel: int64(result.AccessLevel),
	}, nil
}

func (c *AuthController) BanAccount(ctx context.Context, request *authpb.BanAccountRequest) (*authpb.Empty, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthController.BanAccount")
	defer span.End()

	if request.GetSubject() == 0 {
		return nil, status.Error(codes.InvalidArgument, "subject is required")
	}

	err := c.authService.BanAccount(
		ctx,
		int(request.GetSubject()),
		int(request.GetBanUntilUnix()),
		request.GetReason(),
	)
	if err != nil {
		return nil, err
	}

	return &authpb.Empty{}, nil
}
