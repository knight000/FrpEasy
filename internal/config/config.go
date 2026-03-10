package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	Presets []PresetConfig `toml:"presets"`
}

type PresetConfig struct {
	ID       string          `toml:"id"`
	Name     string          `toml:"name"`
	Servers  []ServerConfig  `toml:"servers"`
	Services []ServiceConfig `toml:"services"`
}

type ServerConfig struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Token   string `toml:"token"`
	Enabled bool   `toml:"enabled"`
}

type ServiceConfig struct {
	ID             string `toml:"id"`
	Name           string `toml:"name"`
	Protocol       string `toml:"protocol"`
	LocalIP        string `toml:"local_ip"`
	LocalPort      int    `toml:"local_port"`
	RemotePort     int    `toml:"remote_port"`
	UseEncryption  bool   `toml:"use_encryption"`
	UseCompression bool   `toml:"use_compression"`
}

func SaveConfig(path string, config *AppConfig) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config AppConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func ParseConfigString(tomlStr string) (*AppConfig, error) {
	var config AppConfig
	if err := toml.Unmarshal([]byte(tomlStr), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &config, nil
}

func ToTomlString(config *AppConfig) (string, error) {
	data, err := toml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}
	return string(data), nil
}
