package access_level

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

//go:generate stringer -type=AccessLevel
type AccessLevel int

const (
	User AccessLevel = iota
	ViewAllUsers
	ViewInventory
	ViewMatches
	Admin
	CreateItem
	GiveItem
	RevokeItem
	UpdateItem
	ResetHwid
	AddAdmin
	DeleteItem
	Dev
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
		for i := User; i <= Dev; i++ {
			if strings.EqualFold(i.String(), typedValue) {
				*a = i

				return nil
			}
		}

		return fmt.Errorf("%w: %v", errUnknownAccessLevel, value)
	case []byte:
		return a.Scan(string(typedValue))
	case int64:
		if typedValue >= int64(User) && typedValue <= int64(Dev) {
			*a = AccessLevel(typedValue)

			return nil
		}

		return fmt.Errorf("%w (%v)", errOutOfRangeAccessLevel, value)
	default:
		return fmt.Errorf("%w: %T", errInvalidTypeAccessLevel, value)
	}
}

// FromStringOrDefault converts string to AccessLevel (optional helper).
func FromStringOrDefault(s string) AccessLevel {
	for i := User; i <= Dev; i++ {
		if strings.EqualFold(i.String(), s) {
			return i
		}
	}

	return 0
}
