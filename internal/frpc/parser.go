package frpc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"frpeasy/internal/models"

	"github.com/pelletier/go-toml/v2"
)

type FrpConfig struct {
	ServerAddr     string
	ServerPort     int
	Token          string
	UseEncryption  bool
	UseCompression bool
	Proxies        []ProxyConfig
}

type ProxyConfig struct {
	Name           string
	Type           string
	LocalIP        string
	LocalPort      int
	RemotePort     int
	UseEncryption  bool
	UseCompression bool
}

func ParseTomlFile(filePath string) (*FrpConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseTomlContent(content)
}

func ParseTomlContent(content []byte) (*FrpConfig, error) {
	var raw map[string]interface{}
	if err := toml.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	config := &FrpConfig{
		ServerPort: 7000,
	}

	if v, ok := raw["serverAddr"].(string); ok {
		config.ServerAddr = v
	}
	if v, ok := raw["serverAddr"].(string); ok {
		config.ServerAddr = v
	}
	if v, ok := raw["serverPort"].(int64); ok {
		config.ServerPort = int(v)
	}

	if auth, ok := raw["auth"].(map[string]interface{}); ok {
		if t, ok := auth["token"].(string); ok {
			config.Token = t
		}
	}
	if token, ok := raw["auth.token"].(string); ok {
		config.Token = token
	}

	if transport, ok := raw["transport"].(map[string]interface{}); ok {
		if v, ok := transport["useEncryption"].(bool); ok {
			config.UseEncryption = v
		}
		if v, ok := transport["useCompression"].(bool); ok {
			config.UseCompression = v
		}
	}

	if proxies, ok := raw["proxies"].([]interface{}); ok {
		for _, p := range proxies {
			if proxy, ok := p.(map[string]interface{}); ok {
				proxyConfig := ProxyConfig{
					LocalIP: "127.0.0.1",
				}

				if v, ok := proxy["name"].(string); ok {
					proxyConfig.Name = v
				}
				if v, ok := proxy["type"].(string); ok {
					proxyConfig.Type = v
				}
				if v, ok := proxy["localIP"].(string); ok {
					proxyConfig.LocalIP = v
				}
				if v, ok := proxy["localPort"].(int64); ok {
					proxyConfig.LocalPort = int(v)
				}
				if v, ok := proxy["remotePort"].(int64); ok {
					proxyConfig.RemotePort = int(v)
				}

				if transport, ok := proxy["transport"].(map[string]interface{}); ok {
					if v, ok := transport["useEncryption"].(bool); ok {
						proxyConfig.UseEncryption = v
					}
					if v, ok := transport["useCompression"].(bool); ok {
						proxyConfig.UseCompression = v
					}
				}

				config.Proxies = append(config.Proxies, proxyConfig)
			}
		}
	}

	return config, nil
}

func ParseIniFile(filePath string) (*FrpConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseIniContent(content)
}

func ParseIniContent(content []byte) (*FrpConfig, error) {
	config := &FrpConfig{
		ServerPort: 7000,
	}

	sections := parseIniSections(string(content))

	if common, ok := sections["common"]; ok {
		if v, ok := common["server_addr"]; ok {
			config.ServerAddr = v
		}
		if v, ok := common["serveraddr"]; ok {
			config.ServerAddr = v
		}
		if v, ok := common["server_port"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				config.ServerPort = port
			}
		}
		if v, ok := common["serverport"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				config.ServerPort = port
			}
		}
		if v, ok := common["token"]; ok {
			config.Token = v
		}
		if v, ok := common["auth_token"]; ok {
			config.Token = v
		}
		if v, ok := common["use_encryption"]; ok {
			config.UseEncryption = strings.ToLower(v) == "true"
		}
		if v, ok := common["useencryption"]; ok {
			config.UseEncryption = strings.ToLower(v) == "true"
		}
		if v, ok := common["use_compression"]; ok {
			config.UseCompression = strings.ToLower(v) == "true"
		}
		if v, ok := common["usecompression"]; ok {
			config.UseCompression = strings.ToLower(v) == "true"
		}
	}

	for name, section := range sections {
		if name == "common" || name == "range:" {
			continue
		}

		if strings.HasPrefix(name, "range:") {
			continue
		}

		proxyConfig := ProxyConfig{
			Name:    name,
			Type:    "tcp",
			LocalIP: "127.0.0.1",
		}

		if v, ok := section["type"]; ok {
			proxyConfig.Type = v
		}
		if v, ok := section["local_ip"]; ok {
			proxyConfig.LocalIP = v
		}
		if v, ok := section["localip"]; ok {
			proxyConfig.LocalIP = v
		}
		if v, ok := section["local_port"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				proxyConfig.LocalPort = port
			}
		}
		if v, ok := section["localport"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				proxyConfig.LocalPort = port
			}
		}
		if v, ok := section["remote_port"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				proxyConfig.RemotePort = port
			}
		}
		if v, ok := section["remoteport"]; ok {
			if port, err := strconv.Atoi(v); err == nil {
				proxyConfig.RemotePort = port
			}
		}
		if v, ok := section["use_encryption"]; ok {
			proxyConfig.UseEncryption = strings.ToLower(v) == "true"
		} else {
			proxyConfig.UseEncryption = config.UseEncryption
		}
		if v, ok := section["use_compression"]; ok {
			proxyConfig.UseCompression = strings.ToLower(v) == "true"
		} else {
			proxyConfig.UseCompression = config.UseCompression
		}

		config.Proxies = append(config.Proxies, proxyConfig)
	}

	return config, nil
}

func parseIniSections(content string) map[string]map[string]string {
	sections := make(map[string]map[string]string)
	var currentSection string

	scanner := bufio.NewScanner(strings.NewReader(content))
	sectionRegex := regexp.MustCompile(`^\s*\[([^\]]+)\]\s*$`)
	kvRegex := regexp.MustCompile(`^\s*([a-zA-Z0-9_]+)\s*=\s*(.*)\s*$`)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := sectionRegex.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			sections[currentSection] = make(map[string]string)
			continue
		}

		if currentSection == "" {
			continue
		}

		if matches := kvRegex.FindStringSubmatch(line); matches != nil {
			key := strings.ToLower(matches[1])
			value := strings.Trim(matches[2], "\"'")
			sections[currentSection][key] = value
		}
	}

	return sections
}

func ParseFrpConfigFile(filePath string) (*FrpConfig, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".toml":
		return ParseTomlFile(filePath)
	case ".ini", ".conf":
		return ParseIniFile(filePath)
	default:
		if content, err := os.ReadFile(filePath); err == nil {
			if strings.Contains(string(content), "[[proxies]]") || strings.HasPrefix(strings.TrimSpace(string(content)), "serverAddr") {
				return ParseTomlContent(content)
			}
			return ParseIniContent(content)
		}
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func ConvertToModels(config *FrpConfig, presetName string) (*models.Preset, error) {
	preset := &models.Preset{
		ID:       generateID(),
		Name:     presetName,
		Servers:  []models.Server{},
		Services: []models.Service{},
	}

	server := models.Server{
		ID:      generateID(),
		Name:    "主服务器",
		Address: config.ServerAddr,
		Port:    config.ServerPort,
		Token:   config.Token,
		Status:  models.StatusOffline,
		Enabled: false,
		Logs:    []models.LogEntry{},
		Uptime:  0,
	}

	preset.Servers = append(preset.Servers, server)

	for _, proxy := range config.Proxies {
		useEncryption := proxy.UseEncryption
		useCompression := proxy.UseCompression

		if !useEncryption && config.UseEncryption {
			useEncryption = config.UseEncryption
		}
		if !useCompression && config.UseCompression {
			useCompression = config.UseCompression
		}

		service := models.Service{
			ID:             generateID(),
			Name:           proxy.Name,
			Protocol:       models.ServiceProtocol(strings.ToUpper(proxy.Type)),
			LocalIP:        proxy.LocalIP,
			LocalPort:      proxy.LocalPort,
			RemotePort:     proxy.RemotePort,
			UseEncryption:  useEncryption,
			UseCompression: useCompression,
		}
		preset.Services = append(preset.Services, service)
	}

	return preset, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
}
