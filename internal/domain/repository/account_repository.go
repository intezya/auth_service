package repository

import (
	"context"
	"github.com/intezya/auth_service/internal/domain/dto"
	"time"
)

type AccountRepository interface {
	Create(ctx context.Context, username string, password string, hardwareId string) (*dto.AccountDTO, error)
	FindByID(ctx context.Context, id int) (*dto.AccountDTO, error)
	FindByLowerUsername(ctx context.Context, username string) (*dto.AccountDTO, error)
	UpdateHardwareIDByID(ctx context.Context, id int, hardwareId string) error
	UpdateBannedUntilBannedReasonByID(ctx context.Context, id int, bannedUntil *time.Time, banReason *string) error
}
