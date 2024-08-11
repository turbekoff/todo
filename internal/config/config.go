package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Host           string        `env:"HOST" env-default:"0.0.0.0"`
	Port           int           `env:"PORT" env-default:"80"`
	ReadTimeout    time.Duration `env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout   time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	MaxHeaderBytes int           `env:"MAX_HEADER_BYTES" env-default:"1"`
}

type JWTConfig struct {
	SigningKey      string        `env:"SIGNING_KEY" env-required:"true"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" env-default:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"720h"`
}

type MongoConfig struct {
	URI      string `env:"URI" env-required:"true"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type Config struct {
	DebugMode      Mode        `env:"DEBUG_MODE" env-default:"development"`
	PasswordPepper string      `env:"PASSWORD_PEPPER" env-required:"true"`
	JWT            JWTConfig   `env-prefix:"JWT_"`
	HTTP           HTTPConfig  `env-prefix:"HTTP_"`
	Mongo          MongoConfig `env-prefix:"MONGO_"`
}

func Load() (Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{DebugMode: M_NULL}, err
	}
	return config, nil
}
