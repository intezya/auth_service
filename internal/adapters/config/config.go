package config

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/intezya/auth_service/internal/adapters/grpc"
	"github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
	"github.com/intezya/auth_service/internal/infrastructure/persistence"
	"github.com/intezya/auth_service/internal/pkg/crypto"
	"github.com/intezya/auth_service/internal/pkg/jwt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Logger LoggerConfig
	Server grpc.Config
	Tracer tracer.Config
	JWT    jwt.Config
	Crypto crypto.Config
	Ent    persistence.EntConfig

	EnvType string `env:"ENV" env-default:"dev"` // dev / prod
}

var Cfg Config

func loadEnvFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
}

func LoadConfig() Config {
	envFile := flag.String("env-file", "", "env file path")
	flag.Parse()

	if envFile != nil && *envFile != "" {
		if err := loadEnvFromFile(*envFile); err != nil {
			slog.Warn(fmt.Sprintf("failed to load env file %s: %v", *envFile, err))
		}
	}

	err := cleanenv.ReadEnv(&Cfg)

	if err != nil {
		panic(fmt.Sprintf("failed to load environment variables: %v", err))
	}

	return Cfg
}
