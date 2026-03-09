package config

import (
	"encoding/json"
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

type JSONConfig struct {
	Presets []JSONPreset `json:"presets"`
}

type JSONPreset struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Servers  []JSONServer  `json:"servers"`
	Services []JSONService `json:"services"`
}

type JSONServer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Token   string `json:"token"`
	Enabled bool   `json:"enabled"`
}

type JSONService struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Protocol       string `json:"protocol"`
	LocalIP        string `json:"localIp"`
	LocalPort      int    `json:"localPort"`
	RemotePort     int    `json:"remotePort"`
	UseEncryption  bool   `json:"useEncryption"`
	UseCompression bool   `json:"useCompression"`
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

func ToJSON(config *AppConfig) string {
	var jsonConfig JSONConfig
	for _, p := range config.Presets {
		preset := JSONPreset{
			ID:   p.ID,
			Name: p.Name,
		}
		for _, s := range p.Servers {
			preset.Servers = append(preset.Servers, JSONServer{
				ID:      s.ID,
				Name:    s.Name,
				Address: s.Address,
				Port:    s.Port,
				Token:   s.Token,
				Enabled: s.Enabled,
			})
		}
		for _, s := range p.Services {
			preset.Services = append(preset.Services, JSONService{
				ID:             s.ID,
				Name:           s.Name,
				Protocol:       s.Protocol,
				LocalIP:        s.LocalIP,
				LocalPort:      s.LocalPort,
				RemotePort:     s.RemotePort,
				UseEncryption:  s.UseEncryption,
				UseCompression: s.UseCompression,
			})
		}
		jsonConfig.Presets = append(jsonConfig.Presets, preset)
	}

	data, _ := json.Marshal(jsonConfig)
	return string(data)
}

func FromJSON(jsonStr string) (*AppConfig, error) {
	var jsonConfig JSONConfig
	if err := json.Unmarshal([]byte(jsonStr), &jsonConfig); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var config AppConfig
	for _, p := range jsonConfig.Presets {
		preset := PresetConfig{
			ID:   p.ID,
			Name: p.Name,
		}
		for _, s := range p.Servers {
			preset.Servers = append(preset.Servers, ServerConfig{
				ID:      s.ID,
				Name:    s.Name,
				Address: s.Address,
				Port:    s.Port,
				Token:   s.Token,
				Enabled: s.Enabled,
			})
		}
		for _, s := range p.Services {
			preset.Services = append(preset.Services, ServiceConfig{
				ID:             s.ID,
				Name:           s.Name,
				Protocol:       s.Protocol,
				LocalIP:        s.LocalIP,
				LocalPort:      s.LocalPort,
				RemotePort:     s.RemotePort,
				UseEncryption:  s.UseEncryption,
				UseCompression: s.UseCompression,
			})
		}
		config.Presets = append(config.Presets, preset)
	}

	return &config, nil
}
