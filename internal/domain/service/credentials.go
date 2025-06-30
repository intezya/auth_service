package service

import (
	"context"
	"github.com/intezya/auth_service/internal/domain/dto"
)

type PasswordEncoder interface {
	EncodePassword(ctx context.Context, password string) string
	VerifyPassword(ctx context.Context, password, hash string) bool
	EncodeHardwareID(ctx context.Context, hardwareID string) string
	VerifyHardwareID(ctx context.Context, hardwareID, hash string) bool
}

type TokenManager interface {
	Generate(accountID int) string
	Parse(token string) (*dto.TokenData, error)
}
