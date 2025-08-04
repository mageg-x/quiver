package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 全局配置结构体
type Config struct {
	Database map[string]DatabaseConfig `yaml:"database"`
	Server   ServerConfig              `yaml:"server"`
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

var globalConfig *Config

// LoadConfig 加载YAML配置文件
func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	globalConfig = &cfg
	// 缺省值
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}

	return &cfg, nil
}

// GetDatabaseConfig 获取当前环境下的数据库配置
func GetDatabaseConfig(env string) (DatabaseConfig, error) {
	if globalConfig == nil {
		return DatabaseConfig{}, fmt.Errorf("config not loaded")
	}
	dbConf, exists := globalConfig.Database[env]
	if !exists {
		return DatabaseConfig{}, fmt.Errorf("environment %s not found in config", env)
	}
	return dbConf, nil
}

// GetServerConfig 获取服务器配置
func GetServerConfig() ServerConfig {
	if globalConfig == nil {
		return ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		}
	}

	return globalConfig.Server
}
