package persistence

import (
	"context"
	"fmt"
	"github.com/intezya/auth_service/internal/infrastructure/ent"
	"github.com/intezya/auth_service/internal/infrastructure/ent/migrate"

	"time"
)

import _ "github.com/lib/pq" //nolint

const driverName = "postgres" //nolint

const (
	defaultEntReconnectMaxRetries = 5
	defaultEntReconnectDelay      = 2 * time.Second
)

type Logger interface {
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type EntConfig struct {
	source     string
	maxRetries int
	retryDelay time.Duration
	Debug      bool `env:"ORM_DEBUG" env-default:"false"`

	Host     string `env:"DATABASE_HOST" env-default:"localhost"`
	Port     int    `env:"DATABASE_PORT" env-default:"5432"`
	User     string `env:"DATABASE_USERNAME" env-default:"postgres"`
	Password string `env:"DATABASE_PASSWORD" env-default:"postgres"`
	DBName   string `env:"DATABASE_NAME" env-default:"postgres"`
	SSL      string `env:"DATABASE_SSL" env-default:"disable"`
}

func SetupEnt(config EntConfig, logger Logger) *ent.Client {
	maxRetries := gt0(config.maxRetries, defaultEntReconnectMaxRetries)
	retryDelay := gt0(config.retryDelay, defaultEntReconnectDelay)

	source := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSL,
	)

	entClient, err := ent.Open(driverName, source)
	if err != nil {
		logger.Fatal(err) // invalid driver
	}

	if config.Debug {
		entClient = entClient.Debug()
	}

	// Retry connecting to the database if it fails
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = entClient.Schema.Create(
			context.Background(),
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
		)
		if err == nil {
			logger.Infof("Database migrations runned success on attempt %d", attempt)

			break
		}

		logger.Warnf(
			"Attempt %d of %d: Failed to run migrations for database: %v",
			attempt,
			maxRetries,
			err,
		)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logger.Fatalf("failed to create schema (all attempts are over)")
	}

	return entClient
}

func gt0[T int | time.Duration](value T, fallback T) T {
	if value <= 0 {
		return fallback
	}

	return value
}
