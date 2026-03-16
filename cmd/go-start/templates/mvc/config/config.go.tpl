package config

import (
	"fmt"

	"github.com/spf13/viper"
	"{{.Module}}/pkg/database"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig    `yaml:"server" mapstructure:"server"`
	Database database.Config `yaml:"database" mapstructure:"database"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int `yaml:"port" mapstructure:"port"`
}

// Addr 返回监听地址
func (s *ServerConfig) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

// Load 从文件和环境变量加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/{{.ProjectName}}/")
	viper.AddConfigPath("$HOME/.{{.ProjectName}}/")

	viper.SetDefault("server.port", {{.ServerPort}})
	viper.SetDefault("database.driver", "{{.Database}}")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("{{.ProjectName}}")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &cfg, nil
}
