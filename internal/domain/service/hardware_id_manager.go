package service

import (
	"context"
	entity "github.com/intezya/auth_service/internal/domain/account"
	"github.com/intezya/auth_service/internal/domain/repository"
)

type HardwareIDManager interface {
	ValidateAndSetHardwareID(ctx context.Context, account *entity.Account, providedHardwareID string) error
}

type hardwareIDManager struct {
	accountRepository repository.AccountRepository
	passwordEncoder   PasswordEncoder
}

func NewHardwareIDManager(
	accountRepository repository.AccountRepository,
	passwordEncoder PasswordEncoder,
) HardwareIDManager {
	return &hardwareIDManager{
		accountRepository: accountRepository,
		passwordEncoder:   passwordEncoder,
	}
}

func (h *hardwareIDManager) ValidateAndSetHardwareID(
	ctx context.Context,
	account *entity.Account,
	providedHardwareID string,
) error {
	if account.HardwareID() != nil {
		if h.passwordEncoder.VerifyHardwareID(ctx, providedHardwareID, *account.HardwareID()) {
			panic("TODO()")
			//return TODO()
		}
		return nil
	} else {
		hashedHardwareID := h.passwordEncoder.EncodeHardwareID(ctx, providedHardwareID)
		account.SetHardwareID(entity.HardwareID(hashedHardwareID))
		if err := h.accountRepository.Update(ctx, account); err != nil {
			return err // hardware id conflict
		}
	}

	return nil
}
