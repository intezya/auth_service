package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/intezya/auth_service/internal/domain/dto"
	"github.com/intezya/pkglib/crypto"
	"strconv"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type Config struct {
	SecretKey      string        `env:"JWT_SECRET_KEY" env-required:"true"`
	Issuer         string        `env:"JWT_ISSUER" env-required:"true"`
	ExpirationTime time.Duration `env:"JWT_EXPIRATION_TIME" env-default:"24h"`
}

type Claim struct {
	jwt.RegisteredClaims
}

type TokenHelper struct {
	secretKey      []byte
	issuer         string
	expirationTime time.Duration
}

func NewTokenHelper(config Config) *TokenHelper {
	secretKey := []byte(config.SecretKey)

	if len(config.SecretKey) < 32 {
		secretKey = []byte(crypto.HashSHA256(config.SecretKey))
	}

	return &TokenHelper{
		secretKey:      secretKey,
		issuer:         config.Issuer,
		expirationTime: config.ExpirationTime,
	}
}

func (t *TokenHelper) Generate(subject int) string {
	claims := &Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   strconv.Itoa(subject),
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.expirationTime)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(t.secretKey)

	return tokenString
}

func (t *TokenHelper) Parse(tokenString string) (*dto.DataFromToken, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return t.secretKey, nil
		},
		jwt.WithIssuer(t.issuer),
		jwt.WithStrictDecoding(),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	subj, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, err
	}

	return &dto.DataFromToken{
		ID:        claims.ID,
		Subject:   subj,
		Issuer:    claims.Issuer,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
