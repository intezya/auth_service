package crypto

import "errors"

var errInvalidEncodeFormat = errors.New("invalid encode format")

type Config struct {
	HardwareIDEncryptionKey string `env:"HARDWARE_ID_ENCRYPTION_KEY" env-required:"true"`
}
