package mapper

import (
	"github.com/intezya/auth_service/internal/domain/dto"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
)

func AccountToDto(account *ent.Account) *dto.AccountDTO {
	return &dto.AccountDTO{
		ID:          account.ID,
		Username:    account.Username,
		Password:    account.Password,
		HardwareID:  account.HardwareID,
		AccessLevel: int(account.AccessLevel),
		CreatedAt:   account.CreatedAt,
		BannedUntil: account.BannedUntil,
		BanReason:   account.BanReason,
	}
}
