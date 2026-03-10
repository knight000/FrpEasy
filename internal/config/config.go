package config

import (
	"fmt"
	"os"

	"frpeasy/internal/models"

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
	AdvancedConfig string `toml:"advanced_config"`
	IsAdvanced     bool   `toml:"is_advanced"`
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

func ServerConfigFromModel(s models.Server) ServerConfig {
	return ServerConfig{
		ID:      s.ID,
		Name:    s.Name,
		Address: s.Address,
		Port:    s.Port,
		Token:   s.Token,
		Enabled: s.Enabled,
	}
}

func ServiceConfigFromModel(s models.Service) ServiceConfig {
	return ServiceConfig{
		ID:             s.ID,
		Name:           s.Name,
		Protocol:       string(s.Protocol),
		LocalIP:        s.LocalIP,
		LocalPort:      s.LocalPort,
		RemotePort:     s.RemotePort,
		UseEncryption:  s.UseEncryption,
		UseCompression: s.UseCompression,
		AdvancedConfig: s.AdvancedConfig,
		IsAdvanced:     s.IsAdvanced,
	}
}

func PresetConfigFromModel(p models.Preset) PresetConfig {
	servers := make([]ServerConfig, len(p.Servers))
	for i, s := range p.Servers {
		servers[i] = ServerConfigFromModel(s)
	}
	services := make([]ServiceConfig, len(p.Services))
	for i, s := range p.Services {
		services[i] = ServiceConfigFromModel(s)
	}
	return PresetConfig{
		ID:       p.ID,
		Name:     p.Name,
		Servers:  servers,
		Services: services,
	}
}

func (s ServerConfig) ToModel() models.Server {
	return models.Server{
		ID:      s.ID,
		Name:    s.Name,
		Address: s.Address,
		Port:    s.Port,
		Token:   s.Token,
		Enabled: s.Enabled,
		Status:  models.StatusOffline,
	}
}

func (s ServiceConfig) ToModel() models.Service {
	return models.Service{
		ID:             s.ID,
		Name:           s.Name,
		Protocol:       models.ServiceProtocol(s.Protocol),
		LocalIP:        s.LocalIP,
		LocalPort:      s.LocalPort,
		RemotePort:     s.RemotePort,
		UseEncryption:  s.UseEncryption,
		UseCompression: s.UseCompression,
		AdvancedConfig: s.AdvancedConfig,
		IsAdvanced:     s.IsAdvanced,
	}
}

func (p PresetConfig) ToModel() models.Preset {
	servers := make([]models.Server, len(p.Servers))
	for i, s := range p.Servers {
		servers[i] = s.ToModel()
	}
	services := make([]models.Service, len(p.Services))
	for i, s := range p.Services {
		services[i] = s.ToModel()
	}
	return models.Preset{
		ID:       p.ID,
		Name:     p.Name,
		Servers:  servers,
		Services: services,
	}
}
