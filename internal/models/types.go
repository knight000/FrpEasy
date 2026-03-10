package models

type ServerStatus string

const (
	StatusOnline     ServerStatus = "online"
	StatusOffline    ServerStatus = "offline"
	StatusConnecting ServerStatus = "connecting"
	StatusError      ServerStatus = "error"
)

type ServiceProtocol string

const (
	ProtocolTCP   ServiceProtocol = "TCP"
	ProtocolUDP   ServiceProtocol = "UDP"
	ProtocolHTTP  ServiceProtocol = "HTTP"
	ProtocolHTTPS ServiceProtocol = "HTTPS"
)

type LogEntry struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

type Service struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Protocol       ServiceProtocol `json:"protocol"`
	LocalIP        string          `json:"local_ip"`
	LocalPort      int             `json:"local_port"`
	RemotePort     int             `json:"remote_port"`
	UseEncryption  bool            `json:"use_encryption"`
	UseCompression bool            `json:"use_compression"`
	AdvancedConfig string          `json:"advanced_config"`
	IsAdvanced     bool            `json:"is_advanced"`
}

type Server struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Port    int          `json:"port"`
	Token   string       `json:"token"`
	Status  ServerStatus `json:"status"`
	Enabled bool         `json:"enabled"`
	Logs    []LogEntry   `json:"logs"`
	Uptime  int          `json:"uptime"`
}

type Preset struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Servers  []Server  `json:"servers"`
	Services []Service `json:"services"`
}

type DownloadProgress struct {
	TotalBytes        int64   `json:"total_bytes"`
	Downloaded        int64   `json:"downloaded"`
	Percentage        float64 `json:"percentage"`
	IsComplete        bool    `json:"is_complete"`
	IsError           bool    `json:"is_error"`
	ErrorMessage      string  `json:"error_message,omitempty"`
	DownloadedVersion string  `json:"downloaded_version,omitempty"`
	VersionFetchError string  `json:"version_fetch_error,omitempty"`
}

type ServerRuntime struct {
	PresetID   string `json:"preset_id"`
	ServerID   string `json:"server_id"`
	ProcessPID int    `json:"process_pid"`
	ConfigPath string `json:"config_path"`
}
