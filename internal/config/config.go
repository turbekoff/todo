package config

import (
	"errors"
	"os"
	"path/filepath"
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

func loadPath() (string, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return path, nil
	}

	root, err := os.Getwd()
	if err != nil {
		return "", errors.New("root path not found")
	}
	return filepath.Join(root, path), nil
}

func Load() (Config, error) {
	path, err := loadPath()
	if err != nil {
		return Config{DebugMode: M_NULL}, err
	}

	var config Config
	if path != "" {
		err = cleanenv.ReadConfig(path, &config)
	} else {
		err = cleanenv.ReadEnv(&config)
	}

	if err != nil {
		return Config{DebugMode: M_NULL}, err
	}
	return config, nil
}
