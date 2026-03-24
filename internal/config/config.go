package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	App    AppConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            string
	ShutdownTimeout int // seconds
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

// AppConfig holds application-specific configuration.
type AppConfig struct {
	BaseURL      string
	LogLevel     string
	MinCodeLen   int // Minimum length for short codes
	DefaultTTL   int // Default TTL in seconds (0 = no expiry)
	MaxURLLength int // Maximum allowed URL length
}

// Load loads configuration from environment variables with sensible defaults.
// TODO: Implement configuration loading from environment variables.
// TODO: Add validation for required fields.
// TODO: Consider adding support for config files (YAML/JSON).
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 30),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 50),
		},
		App: AppConfig{
			BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
			LogLevel:     getEnv("LOG_LEVEL", "info"),
			MinCodeLen:   getEnvAsInt("MIN_CODE_LENGTH", 6),
			DefaultTTL:   getEnvAsInt("DEFAULT_TTL", 0),
			MaxURLLength: getEnvAsInt("MAX_URL_LENGTH", 2048),
		},
	}

	// TODO: Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration.
// TODO: Implement comprehensive validation logic.
func (c *Config) Validate() error {
	// TODO: Validate server port format
	// TODO: Validate Redis address format
	// TODO: Validate base URL format
	// TODO: Validate numeric ranges (pool size, timeouts, etc.)
	return nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
