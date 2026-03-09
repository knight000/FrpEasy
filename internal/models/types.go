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
	LocalIP        string          `json:"localIp"`
	LocalPort      int             `json:"localPort"`
	RemotePort     int             `json:"remotePort"`
	UseEncryption  bool            `json:"useEncryption"`
	UseCompression bool            `json:"useCompression"`
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
	TotalBytes   int64   `json:"totalBytes"`
	Downloaded   int64   `json:"downloaded"`
	Percentage   float64 `json:"percentage"`
	IsComplete   bool    `json:"isComplete"`
	IsError      bool    `json:"isError"`
	ErrorMessage string  `json:"errorMessage"`
}

type ServerRuntime struct {
	PresetID   string `json:"presetId"`
	ServerID   string `json:"serverId"`
	ProcessPID int    `json:"processPid"`
	ConfigPath string `json:"configPath"`
}
