package dto

import "time"

type TokenData struct {
	ID          string    `json:"id"`
	AccessLevel int       `json:"access_level"`
	Subject     int       `json:"subject"`
	Issuer      string    `json:"issuer"`
	IssuedAt    time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
