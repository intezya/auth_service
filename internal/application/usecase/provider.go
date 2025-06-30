package usecase

import (
	"github.com/intezya/auth_service/internal/domain/service"
	"github.com/intezya/auth_service/internal/infrastructure/persistence"
	"github.com/intezya/auth_service/internal/pkg/validator"
)

type Provider struct {
	AuthUseCase AuthUseCase
}

func NewProvider(
	repositoryProvider *persistence.Provider,
	validatorProvider *domainvalidator.Provider,
	passwordEncoder service.PasswordEncoder,
	tokenManager service.TokenManager,
	hardwareIDManager service.HardwareIDManager,
) *Provider {
	return &Provider{
		AuthUseCase: NewAuthUseCase(
			repositoryProvider.AccountRepository,
			passwordEncoder,
			tokenManager,
			hardwareIDManager,
			validatorProvider.UsernameValidator,
			validatorProvider.PasswordValidator,
			validatorProvider.HardwareValidator,
		),
	}
}
