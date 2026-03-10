package frpc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"frpeasy/internal/models"

	"github.com/google/uuid"
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
	AdvancedConfig string
	IsAdvanced     bool
}

func ParseTomlFile(filePath string) (*FrpConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseTomlContent(content)
}

func isBasicProxyField(key string) bool {
	lowerKey := strings.ToLower(key)
	switch lowerKey {
	case "name", "type", "localip", "localport", "remoteport",
		"transport.useencryption", "transport.usecompression":
		return true
	case "transport":
		return false
	}
	return false
}

func getMapStringCI(m map[string]interface{}, key string) string {
	lowerKey := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			if s, ok := v.(string); ok {
				return s
			}
			break
		}
	}
	return ""
}

func getMapIntCI(m map[string]interface{}, key string) int64 {
	lowerKey := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			if i, ok := v.(int64); ok {
				return i
			}
			break
		}
	}
	return 0
}

func getMapBoolCI(m map[string]interface{}, key string) bool {
	lowerKey := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			if b, ok := v.(bool); ok {
				return b
			}
			break
		}
	}
	return false
}

func getMapMapCI(m map[string]interface{}, key string) map[string]interface{} {
	lowerKey := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			if sub, ok := v.(map[string]interface{}); ok {
				return sub
			}
			break
		}
	}
	return nil
}

func extractProxyToml(proxy map[string]interface{}) (string, bool) {
	var sb strings.Builder
	hasAdvanced := false

	for key, value := range proxy {
		if strings.EqualFold(key, "transport") {
			if transport, ok := value.(map[string]interface{}); ok {
				for tkey, tvalue := range transport {
					fullKey := "transport." + tkey
					if isBasicProxyField(fullKey) {
						if b, ok := tvalue.(bool); ok && b {
							sb.WriteString(fmt.Sprintf("%s = true\n", fullKey))
						}
					} else {
						hasAdvanced = true
						switch v := tvalue.(type) {
						case string:
							sb.WriteString(fmt.Sprintf("%s = \"%s\"\n", fullKey, v))
						case int64:
							sb.WriteString(fmt.Sprintf("%s = %d\n", fullKey, v))
						case float64:
							sb.WriteString(fmt.Sprintf("%s = %v\n", fullKey, v))
						case bool:
							sb.WriteString(fmt.Sprintf("%s = %v\n", fullKey, v))
						default:
							data, _ := toml.Marshal(map[string]interface{}{tkey: v})
							sb.WriteString(string(data))
						}
					}
				}
			}
		} else if !isBasicProxyField(key) {
			hasAdvanced = true
			switch v := value.(type) {
			case string:
				sb.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, v))
			case int64:
				sb.WriteString(fmt.Sprintf("%s = %d\n", key, v))
			case float64:
				sb.WriteString(fmt.Sprintf("%s = %v\n", key, v))
			case bool:
				sb.WriteString(fmt.Sprintf("%s = %v\n", key, v))
			default:
				data, _ := toml.Marshal(map[string]interface{}{key: v})
				sb.WriteString(string(data))
			}
		}
	}

	return sb.String(), hasAdvanced
}

func parseTomlWithGoTemplate(content []byte) (*FrpConfig, error) {
	config := &FrpConfig{
		ServerPort: 7000,
	}

	contentStr := string(content)

	var raw map[string]interface{}
	toml.Unmarshal(content, &raw)

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

	templateBlocks := extractGoTemplateBlocks(contentStr)
	for _, block := range templateBlocks {
		displayInfo, err := ParseGoTemplateBlock(block)
		if err != nil {
			continue
		}

		frpeasyPrefix := fmt.Sprintf("#FRPEASY#name=%s#protocol=%s#ports=%s\n",
			displayInfo.NamePattern,
			displayInfo.Protocol,
			displayInfo.RemotePorts)

		proxyConfig := ProxyConfig{
			Name:           displayInfo.NamePattern,
			Type:           strings.ToLower(displayInfo.Protocol),
			LocalIP:        "127.0.0.1",
			AdvancedConfig: frpeasyPrefix + block,
			IsAdvanced:     true,
		}

		config.Proxies = append(config.Proxies, proxyConfig)
	}

	return config, nil
}

func extractGoTemplateBlocks(content string) []string {
	var blocks []string

	re := regexp.MustCompile(`(?s)(\{\{-?\s*range[^}]+\}\}(?:.+?)\{\{-?\s*end\s*\}\})`)
	matches := re.FindAllString(content, -1)

	for _, match := range matches {
		if strings.Contains(match, "[[proxies]]") {
			blocks = append(blocks, match)
		}
	}

	return blocks
}

func ParseTomlContent(content []byte) (*FrpConfig, error) {
	contentStr := string(content)

	if ContainsGoTemplate(contentStr) {
		return parseTomlWithGoTemplate(content)
	}

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

				if v := getMapStringCI(proxy, "name"); v != "" {
					proxyConfig.Name = v
				}
				if v := getMapStringCI(proxy, "type"); v != "" {
					proxyConfig.Type = v
				}
				if v := getMapStringCI(proxy, "localIP"); v != "" {
					proxyConfig.LocalIP = v
				}
				if v := getMapIntCI(proxy, "localPort"); v != 0 {
					proxyConfig.LocalPort = int(v)
				}
				if v := getMapIntCI(proxy, "remotePort"); v != 0 {
					proxyConfig.RemotePort = int(v)
				}

				if transport := getMapMapCI(proxy, "transport"); transport != nil {
					if v := getMapBoolCI(transport, "useEncryption"); v {
						proxyConfig.UseEncryption = v
					}
					if v := getMapBoolCI(transport, "useCompression"); v {
						proxyConfig.UseCompression = v
					}
				}

				advancedToml, hasAdvanced := extractProxyToml(proxy)
				if hasAdvanced {
					proxyConfig.AdvancedConfig = advancedToml
					proxyConfig.IsAdvanced = true
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
			AdvancedConfig: proxy.AdvancedConfig,
			IsAdvanced:     proxy.IsAdvanced,
		}
		preset.Services = append(preset.Services, service)
	}

	return preset, nil
}

func generateID() string {
	return uuid.New().String()[:8]
}
