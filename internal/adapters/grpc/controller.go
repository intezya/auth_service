package grpc

import (
	"context"
	"github.com/intezya/auth_service/internal/application/service"
	authpb "github.com/intezya/auth_service/protos/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authController struct {
	authpb.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAuthController(authService service.AuthService) authpb.AuthServiceServer {
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

	err := c.authService.Register(ctx, request.Username, request.Password, request.HardwareId)
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

	result, err := c.authService.Login(ctx, request.Username, request.Password, request.HardwareId)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResponse{
		Token:       result.Token,
		AccessLevel: int64(result.AccessLevel),
	}, nil
}

func (c *authController) VerifyToken(
	ctx context.Context,
	request *authpb.VerifyTokenRequest,
) (*authpb.VerifyTokenResponse, error) {
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

func (c *authController) BanAccount(ctx context.Context, request *authpb.BanAccountRequest) (*authpb.Empty, error) {
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
