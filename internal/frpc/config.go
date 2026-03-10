package frpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"frpeasy/internal/models"
)

func GenerateConfig(server *models.Server, services []models.Service) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", server.Address))
	sb.WriteString(fmt.Sprintf("serverPort = %d\n", server.Port))

	if server.Token != "" {
		sb.WriteString("\n[auth]\n")
		sb.WriteString(fmt.Sprintf("token = \"%s\"\n", server.Token))
	}

	sb.WriteString("\n[log]\n")
	sb.WriteString("to = \"console\"\n")
	sb.WriteString("level = \"info\"\n")

	for _, service := range services {
		if service.IsAdvanced && service.AdvancedConfig != "" {
			sb.WriteString("\n")
			sb.WriteString(service.AdvancedConfig)
			if !strings.HasSuffix(service.AdvancedConfig, "\n") {
				sb.WriteString("\n")
			}
		} else {
			sb.WriteString("\n[[proxies]]\n")
			sb.WriteString(fmt.Sprintf("name = \"%s\"\n", service.Name))
			sb.WriteString(fmt.Sprintf("type = \"%s\"\n", strings.ToLower(string(service.Protocol))))
			sb.WriteString(fmt.Sprintf("localIP = \"%s\"\n", service.LocalIP))
			sb.WriteString(fmt.Sprintf("localPort = %d\n", service.LocalPort))
			sb.WriteString(fmt.Sprintf("remotePort = %d\n", service.RemotePort))

			if service.UseEncryption {
				sb.WriteString("transport.useEncryption = true\n")
			}
			if service.UseCompression {
				sb.WriteString("transport.useCompression = true\n")
			}
		}
	}

	return sb.String()
}

func WriteConfigFile(configDir, presetID, serverID string, server *models.Server, services []models.Service) (string, error) {
	configContent := GenerateConfig(server, services)

	configPath := filepath.Join(configDir, fmt.Sprintf("%s_%s.toml", presetID, serverID))

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}

	return configPath, nil
}

func DeleteConfigFile(configPath string) error {
	if configPath == "" {
		return nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(configPath)
}
