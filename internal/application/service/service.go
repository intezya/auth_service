package service

import (
	"github.com/intezya/auth_service/internal/domain/dto"
)

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
