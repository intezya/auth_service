package domainvalidator

import "github.com/intezya/auth_service/internal/domain/service"

type Provider struct {
	UsernameValidator service.Validator[string]
	PasswordValidator service.Validator[string]
	HardwareValidator service.Validator[string]
}

func NewProvider() *Provider {
	return &Provider{}
}
