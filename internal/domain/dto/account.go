package dto

import (
	"time"
)

type AccountDTO struct {
	ID          int        `json:"id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Password    string     `json:"-"`
	HardwareID  *string    `json:"-"`
	AccessLevel int        `json:"access_level,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	BannedUntil *time.Time `json:"banned_until,omitempty"`
	BanReason   *string    `json:"ban_reason,omitempty"`
}
