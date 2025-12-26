package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/{{.Module}}/pkg/cache"
	"github.com/{{.Module}}/pkg/database"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database database.Config `yaml:"database" mapstructure:"database"`
	Redis    cache.Config    `yaml:"redis" mapstructure:"redis"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int `yaml:"port" mapstructure:"port"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/{{.ProjectName}}/")
	viper.AddConfigPath("$HOME/.{{.ProjectName}}/")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.dial_timeout", 5)
	viper.SetDefault("redis.read_timeout", 3)
	viper.SetDefault("redis.write_timeout", 3)
	viper.SetDefault("redis.pool_timeout", 4)

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("{{.ProjectName}}")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
