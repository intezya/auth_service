package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type AccountID int
type Username string
type HashedPassword string
type HardwareID string

//go:generate stringer -type=AccessLevel
type AccessLevel int

const (
	AccessLevelUser AccessLevel = iota
	AccessLevelViewAllUsers
	AccessLevelViewInventory
	AccessLevelViewMatches
	AccessLevelAdmin
	AccessLevelCreateItem
	AccessLevelGiveItem
	AccessLevelRevokeItem
	AccessLevelUpdateItem
	AccessLevelResetHwid
	AccessLevelAddAdmin
	AccessLevelDeleteItem
	AccessLevelDev
)

var (
	errUnknownAccessLevel     = errors.New("unknown AccessLevel")
	errOutOfRangeAccessLevel  = errors.New("AccessLevel out of range")
	errInvalidTypeAccessLevel = errors.New("invalid AccessLevel value type")
)

// Value implements the driver.Valuer interface for saving to DB (as string).
func (a AccessLevel) Value() (driver.Value, error) {
	return a.String(), nil
}

// Scan implements the sql.Scanner interface for reading from DB (as string).
func (a *AccessLevel) Scan(value interface{}) error {
	switch typedValue := value.(type) {
	case string:
		for i := AccessLevelUser; i <= AccessLevelDev; i++ {
			if strings.EqualFold(i.String(), typedValue) {
				*a = i

				return nil
			}
		}

		return fmt.Errorf("%w: %v", errUnknownAccessLevel, value)
	case []byte:
		return a.Scan(string(typedValue))
	case int64:
		if typedValue >= int64(AccessLevelUser) && typedValue <= int64(AccessLevelDev) {
			*a = AccessLevel(typedValue)

			return nil
		}

		return fmt.Errorf("%w (%v)", errOutOfRangeAccessLevel, value)
	default:
		return fmt.Errorf("%w: %T", errInvalidTypeAccessLevel, value)
	}
}
