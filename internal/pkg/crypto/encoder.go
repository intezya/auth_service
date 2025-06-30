package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/intezya/auth_service/internal/domain/service"
	"github.com/intezya/pkglib/crypto"
	"github.com/intezya/pkglib/generate"
	"strings"
)

type passwordEncoder struct {
	block cipher.Block
}

func NewPasswordEncoder(config Config) service.PasswordEncoder {
	block, err := aes.NewCipher([]byte(config.HardwareIDEncryptionKey)[:])
	if err != nil {
		panic(err) // not 16/24/32 bytes
	}

	return &passwordEncoder{
		block: block,
	}
}

func (p *passwordEncoder) EncodePassword(ctx context.Context, password string) string {
	hashed, err := crypto.HashArgon2(password, crypto.DefaultArgonParams)
	if err != nil {
		panic(err)
	}

	return hashed
}

func (p *passwordEncoder) VerifyPassword(ctx context.Context, password, hash string) bool {
	ok, err := crypto.VerifyArgon2(hash, password)
	if err != nil {
		return false // malformed
	}

	return ok
}

func (p *passwordEncoder) EncodeHardwareID(ctx context.Context, hardwareID string) string {
	salt := generate.RandomBytes(12) //nolint:mnd

	aesgcm, err := cipher.NewGCM(p.block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, salt, []byte(hardwareID), nil)

	return fmt.Sprintf(
		"%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(ciphertext),
	)
}

func (p *passwordEncoder) VerifyHardwareID(ctx context.Context, hardwareID, hash string) bool {
	decoded, err := p.decodeHardwareID(ctx, hash)
	if err != nil {
		return false // malformed
	}

	return decoded == hardwareID
}

func (p *passwordEncoder) decodeHardwareID(ctx context.Context, hardwareID string) (string, error) {
	parts := strings.Split(hardwareID, ":")
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

	aesgcm, err := cipher.NewGCM(p.block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
