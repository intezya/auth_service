package service

import (
	"context"
	"github.com/intezya/auth_service/internal/domain/dto"
)

type AuthService interface {
	Register(
		ctx context.Context,
		username string,
		password string,
		hardwareId string,
	) error
	Login(
		ctx context.Context,
		username string,
		password string,
		hardwareId string,
	) (*dto.AuthenticationResult, error)
	VerifyToken(ctx context.Context, token string) (*dto.DataFromToken, error)
	BanAccount(ctx context.Context, subject int, banUntilUnix int, reason string) error
}

type TokenHelper interface {
	Generate(subject int) string
	Parse(token string) (*dto.DataFromToken, error)
}

type CredentialsHelper interface {
	EncodePassword(raw string) string
	VerifyPassword(raw, hash string) bool
	EncodeHardwareID(raw string) string
	DecodeHardwareID(encoded string) (string, error)
	VerifyHardwareID(raw, encoded string) bool
}
