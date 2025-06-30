package domain

import (
	"errors"
	"github.com/intezya/auth_service/pkg/clock"
	"time"
)

type Account struct {
	id          AccountID
	username    Username
	password    HashedPassword
	hardwareID  *HardwareID
	accessLevel AccessLevel
	bannedUntil *time.Time
	banReason   *string
	createdAt   time.Time
}

func NewAccount(
	username Username,
	password HashedPassword,
	hardwareID HardwareID,
	clock clock.Clock,
) *Account {
	return &Account{
		username:    username,
		password:    password,
		hardwareID:  &hardwareID,
		accessLevel: AccessLevelUser,
		bannedUntil: nil,
		banReason:   nil,
		createdAt:   clock.Now(),
	}
}

func NewAccountFromRepository(
	id AccountID,
	username Username,
	password HashedPassword,
	hardwareID *HardwareID,
	accessLevel AccessLevel,
	bannedUntil *time.Time,
	banReason *string,
	createdAt time.Time,
) *Account {
	return &Account{
		id:          id,
		username:    username,
		password:    password,
		hardwareID:  hardwareID,
		accessLevel: accessLevel,
		bannedUntil: bannedUntil,
		banReason:   banReason,
		createdAt:   createdAt,
	}
}

func (a *Account) ID() int                 { return int(a.id) }
func (a *Account) Username() string        { return string(a.username) }
func (a *Account) Password() string        { return string(a.password) }
func (a *Account) HardwareID() *string     { return (*string)(a.hardwareID) }
func (a *Account) AccessLevel() int        { return int(a.accessLevel) }
func (a *Account) BannedUntil() *time.Time { return a.bannedUntil }
func (a *Account) BanReason() *string      { return a.banReason }
func (a *Account) CreatedAt() time.Time    { return a.createdAt }

func (a *Account) SetHardwareID(hardwareID HardwareID) {
	a.hardwareID = &hardwareID
}

func (a *Account) Ban(until time.Time, reason *string) error {
	if until.Before(time.Now()) {
		return errors.New("ban time must be in the future") // TODO: update error
	}
	a.bannedUntil = &until
	a.banReason = reason
	return nil
}

func (a *Account) Unban() {
	a.bannedUntil = nil
	a.banReason = nil
}

func (a *Account) IsBanned(clock clock.Clock) bool {
	if a.bannedUntil == nil {
		return false
	}

	return a.bannedUntil.Unix() > clock.Now().Unix()
}
