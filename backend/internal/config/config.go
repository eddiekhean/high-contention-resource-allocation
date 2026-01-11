package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MazeServiceConfig struct {
	URL string `yaml:"url" json:"url"`
}

type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins" json:"allowed_origins"`
}

type Config struct {
	Log         Log               `yaml:"log" json:"log"`
	RateLimit   RateLimit         `yaml:"rate_limit" json:"rate_limit"`
	RedisConfig RedisConfig       `yaml:"redis" json:"redis"`
	MazeService MazeServiceConfig `yaml:"maze_service" json:"maze_service"`
	Cors        CorsConfig        `yaml:"cors" json:"cors"`
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
	// Enum validation - Log level
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	validLogLevel := false
	for _, level := range validLogLevels {
		if config.Log.Level == level {
			validLogLevel = true
			break
		}
	}
	if !validLogLevel {
		return fmt.Errorf("invalid log.level: %s (valid values: %s)", config.Log.Level, strings.Join(validLogLevels, ", "))
	}

	// Set default values for log configuration
	if config.Log.Output == "" {
		config.Log.Output = "stdout" // Default to stdout
	}
	if config.Log.Format == "" {
		config.Log.Format = "json" // Default to json
	}

	// Enum validation - Log format
	validLogFormats := []string{"json", "text"}
	validLogFormat := false
	for _, format := range validLogFormats {
		if config.Log.Format == format {
			validLogFormat = true
			break
		}
	}
	if !validLogFormat {
		return fmt.Errorf("invalid log.format: %s (valid values: %s)", config.Log.Format, strings.Join(validLogFormats, ", "))
	}

	// Enum validation - Log output
	validLogOutputs := []string{"stdout", "file", "both"}
	validLogOutput := false
	for _, output := range validLogOutputs {
		if config.Log.Output == output {
			validLogOutput = true
			break
		}
	}
	if !validLogOutput {
		return fmt.Errorf("invalid log.output: %s (valid values: %s)", config.Log.Output, strings.Join(validLogOutputs, ", "))
	}

	// Validate file path when output includes file
	if config.Log.Output == "file" || config.Log.Output == "both" {
		if config.Log.FilePath == "" {
			return fmt.Errorf("log.file_path is required when log.output is 'file' or 'both'")
		}
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

	// Maze Service
	if mazeURL := os.Getenv("MAZE_SERVICE_URL"); mazeURL != "" {
		cfg.MazeService.URL = mazeURL
	}

	// CORS
	if allowedOrigins := os.Getenv("ALLOWED_ORIGINS"); allowedOrigins != "" {
		cfg.Cors.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}
}
