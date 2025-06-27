package config

type LoggerConfig struct {
	Debug         bool   `env:"DEBUG" env-default:"true"`
	TimeZone      string `env:"LOGGER_TIMEZONE" env-default:"UTC"`
	Environment   string `env:"ENV" env-default:"dev"`
	CallerEnabled bool   `env:"LOGGER_WITH_CALLER" env-default:"true"`
}
