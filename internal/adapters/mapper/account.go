package mapper

import (
	domain "github.com/intezya/auth_service/internal/domain/account"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
)

func EntAccountToDomain(account *ent.Account) *domain.Account {
	return domain.NewAccountFromRepository(
		domain.AccountID(account.ID),
		domain.Username(account.Username),
		domain.HashedPassword(account.Password),
		(*domain.HardwareID)(account.HardwareID),
		account.AccessLevel,
		account.BannedUntil,
		account.BanReason,
		account.CreatedAt,
	)
}
