package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Log struct {
	Level      string `yaml:"level" json:"level" validate:"required,oneof=debug info warn error fatal"` // Required with validation
	Format     string `yaml:"format" json:"format" validate:"omitempty,oneof=json text"`                // Optional: json (default), text
	Output     string `yaml:"output" json:"output" validate:"omitempty,oneof=stdout file both"`         // Optional: stdout (default), file, or both
	FilePath   string `yaml:"file_path" json:"file_path"`                                               // File path when output includes file
	MaxSize    int    `yaml:"max_size" json:"max_size"`                                                 // Maximum file size in MB (default: 100MB)
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`                                           // Maximum number of backup files (default: 5)
	MaxAge     int    `yaml:"max_age" json:"max_age"`                                                   // Maximum age of log files in days (default: 30)
	Compress   bool   `yaml:"compress" json:"compress"`                                                 // Whether to compress old log files (default: true)
}
type RateLimit struct {
	Enabled bool `yaml:"enabled"`
	RPS     int  `yaml:"rps"`   // request / second
	Burst   int  `yaml:"burst"` // cho phép vượt ngắn hạn
}
type S3Config struct {
	Enabled   bool   `yaml:"enabled"`
	AccessKey string `yaml:"access_key" validate:"required_if=Enabled true"`
	SecretKey string `yaml:"secret_key" validate:"required_if=Enabled true"`
	Bucket    string `yaml:"bucket" validate:"required_if=Enabled true"`
	Addr      string `yaml:"addr" validate:"required_if=Enabled true"`
}

type Config struct {
	Log         Log            `yaml:"log" json:"log"`
	RateLimit   RateLimit      `yaml:"rate_limit" json:"rate_limit"`
	RedisConfig RedisConfig    `yaml:"redis" json:"redis"`
	S3          S3Config       `yaml:"s3" json:"s3"`
	Postgres    PostgresConfig `yaml:"postgres" json:"postgres"`
}
type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type PostgresConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Addr     string `yaml:"addr" json:"addr" validate:"required_if=Enabled true"`
	User     string `yaml:"user" json:"user" validate:"required_if=Enabled true"`
	Password string `yaml:"password" json:"password" validate:"required_if=Enabled true"`
	DB       string `yaml:"db" json:"db" validate:"required_if=Enabled true"`
}

// LoadFromFile loads configuration from a specific YAML file
func LoadFromFile(filename string) (*Config, error) {
	config := &Config{}

	// Load from specified YAML file
	if err := loadYAMLFile(config, filename); err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", filename, err)
	}
	overrideFromEnv(config)
	// Validate configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadYAMLFile loads configuration from a specific YAML file
func loadYAMLFile(config *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// validate checks if the configuration is valid
func validate(config *Config) error {
	v := validator.New()

	// Register custom validation if needed, but required_if should work
	if err := v.Struct(config); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logrus.Errorf("Validation error: Field '%s' failed on the '%s' tag", err.Field(), err.Tag())
		}
		return err
	}

	// Set default values for log configuration
	if config.Log.Output == "" {
		config.Log.Output = "stdout" // Default to stdout
	}
	if config.Log.Format == "" {
		config.Log.Format = "json" // Default to json
	}

	// Set default values for log rotation if not specified
	if config.Log.MaxSize <= 0 {
		config.Log.MaxSize = 100 // Default 100MB
	}
	if config.Log.MaxBackups <= 0 {
		config.Log.MaxBackups = 5 // Default 5 backup files
	}
	if config.Log.MaxAge <= 0 {
		config.Log.MaxAge = 30 // Default 30 days
	}

	// Validate file path when output includes file
	if config.Log.Output == "file" || config.Log.Output == "both" {
		if config.Log.FilePath == "" {
			err := fmt.Errorf("log.file_path is required when log.output is 'file' or 'both'")
			logrus.Error(err)
			return err
		}
	}

	return nil
}
func overrideFromEnv(cfg *Config) {
	// Redis
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.RedisConfig.Addr = addr
	}
	if pwd := os.Getenv("REDIS_PASSWORD"); pwd != "" {
		cfg.RedisConfig.Password = pwd
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if v, err := strconv.Atoi(db); err == nil {
			cfg.RedisConfig.DB = v
		}
	}

	// Log
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}

	// ===== S3 =====
	if enabled := os.Getenv("S3_ENABLED"); enabled != "" {
		cfg.S3.Enabled = enabled == "true"
	}
	if key := os.Getenv("S3_ACCESS_KEY"); key != "" {
		cfg.S3.AccessKey = key
	}
	if secret := os.Getenv("S3_SECRET_KEY"); secret != "" {
		cfg.S3.SecretKey = secret
	}
	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		cfg.S3.Bucket = bucket
	}
	if region := os.Getenv("AWS_REGION"); region != "" {
		cfg.S3.Addr = region
	}

	// ===== Postgres =====
	if enabled := os.Getenv("POSTGRES_ENABLED"); enabled != "" {
		cfg.Postgres.Enabled = enabled == "true"
	}
	if addr := os.Getenv("POSTGRES_ADDR"); addr != "" {
		cfg.Postgres.Addr = addr
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.Postgres.User = user
	}
	if pwd := os.Getenv("POSTGRES_PASSWORD"); pwd != "" {
		cfg.Postgres.Password = pwd
	}
	if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
		cfg.Postgres.DB = dbName
	}

}
