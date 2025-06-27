package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/intezya/pkglib/crypto"
	"github.com/intezya/pkglib/generate"
)

var errInvalidEncodeFormat = errors.New("invalid encode format")

type Config struct {
	EncryptionKey string `env:"ENCRYPTION_KEY" env-required:"true"`
}

type HashHelper struct {
	block cipher.Block
}

func NewHashHelper(config Config) *HashHelper {
	key := sha256.Sum256([]byte(config.EncryptionKey))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err) // unreachable, cuz key has sha256 format
	}

	return &HashHelper{
		block: block,
	}
}

func (h *HashHelper) EncodePassword(raw string) string {
	hash, err := crypto.HashArgon2(h.preHash(raw), crypto.DefaultArgonParams)
	if err != nil {
		panic(err)
	}

	return hash
}

func (h *HashHelper) VerifyPassword(raw, hash string) bool {
	ok, _ := crypto.VerifyArgon2(hash, h.preHash(raw))

	return ok
}

func (h *HashHelper) EncodeHardwareID(raw string) string {
	salt := generate.RandomBytes(12) //nolint:mnd

	aesgcm, err := cipher.NewGCM(h.block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, salt, []byte(raw), nil)

	return fmt.Sprintf(
		"%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(ciphertext),
	)
}

func (h *HashHelper) DecodeHardwareID(encoded string) (string, error) {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 { //nolint:mnd
		return "", errInvalidEncodeFormat
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(h.block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (h *HashHelper) VerifyHardwareID(raw, encoded string) bool {
	decoded, err := h.DecodeHardwareID(encoded)
	if err != nil {
		return false
	}

	return decoded == raw
}

func (h *HashHelper) preHash(raw string) string {
	shaSum := sha256.Sum256([]byte(raw))

	return hex.EncodeToString(shaSum[:])
}
