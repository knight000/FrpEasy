package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"frpeasy/internal/config"
	"frpeasy/internal/frpc"
	"frpeasy/internal/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const AppVersion = "1.0.0"
const AppGitHubRepo = "https://github.com/your-repo/frpeasy"

type App struct {
	ctx     context.Context
	dataDir string
	manager *frpc.ProcessManager
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}

	a.dataDir = filepath.Join(filepath.Dir(execPath), "frpeasy")

	binDir := filepath.Join(a.dataDir, "bin")
	configDir := filepath.Join(a.dataDir, "configs")

	if err := os.MkdirAll(a.dataDir, 0755); err != nil {
		fmt.Println("Failed to create data directory:", err)
	}

	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Println("Failed to create bin directory:", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Failed to create configs directory:", err)
	}

	a.manager = frpc.NewProcessManager(binDir, configDir)
}

func (a *App) shutdown(ctx context.Context) {
	fmt.Println("Shutting down FrpEasy...")
	if a.manager != nil {
		a.manager.StopAll()
	}
}

func (a *App) GetDataDir() string {
	return a.dataDir
}

func (a *App) GetFrpcVersion() string {
	version, err := frpc.GetFrpcVersion(filepath.Join(a.dataDir, "bin"))
	if err != nil {
		return ""
	}
	return version
}

func (a *App) IsFrpcDownloaded() bool {
	return frpc.IsFrpcDownloaded(filepath.Join(a.dataDir, "bin"))
}

func (a *App) DownloadFrpc(source string) {
	go func() {
		var downloadSource frpc.DownloadSource
		switch source {
		case "github":
			downloadSource = frpc.SourceGithub
		case "fastgit":
			downloadSource = frpc.SourceFastGit
		case "moeyy":
			downloadSource = frpc.SourceMoeyy
		default:
			downloadSource = frpc.SourceGHProxy
		}

		err := frpc.DownloadFrpc(filepath.Join(a.dataDir, "bin"), downloadSource, func(progress models.DownloadProgress) {
			runtime.EventsEmit(a.ctx, "download:progress", progress)
		})

		if err != nil {
			runtime.EventsEmit(a.ctx, "download:progress", models.DownloadProgress{
				IsError:      true,
				ErrorMessage: err.Error(),
			})
		}
	}()
}

func (a *App) StartServer(presetID, serverID string, server models.Server, services []models.Service) error {
	return a.manager.Start(presetID, serverID, &server, services, func(presetID, serverID string, log models.LogEntry) {
		runtime.EventsEmit(a.ctx, "server:log", map[string]interface{}{
			"preset_id": presetID,
			"server_id": serverID,
			"log":       log,
		})
	})
}

func (a *App) StopServer(presetID, serverID string) error {
	return a.manager.Stop(presetID, serverID)
}

func (a *App) IsServerRunning(presetID, serverID string) bool {
	return a.manager.IsRunning(presetID, serverID)
}

func (a *App) GetServerLogs(presetID, serverID string) []models.LogEntry {
	return a.manager.GetLogs(presetID, serverID)
}

func (a *App) GetServerUptime(presetID, serverID string) int {
	return a.manager.GetUptime(presetID, serverID)
}

func (a *App) GetRunningServers() []models.ServerRuntime {
	return a.manager.GetRunningServers()
}

func (a *App) ExportToml(server models.Server, services []models.Service) string {
	return frpc.GenerateConfig(&server, services)
}

type ImportResult struct {
	Preset *models.Preset `json:"preset"`
	Error  string         `json:"error"`
}

func (a *App) ImportFrpFiles() []ImportResult {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "导入 frp 配置文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "frp 配置文件 (*.toml, *.ini)", Pattern: "*.toml;*.ini;*.conf"},
			{DisplayName: "TOML 文件 (*.toml)", Pattern: "*.toml"},
			{DisplayName: "INI 文件 (*.ini)", Pattern: "*.ini;*.conf"},
		},
	})

	if err != nil {
		return []ImportResult{{Error: err.Error()}}
	}

	if len(files) == 0 {
		return []ImportResult{}
	}

	var results []ImportResult

	for _, file := range files {
		config, err := frpc.ParseFrpConfigFile(file)
		if err != nil {
			results = append(results, ImportResult{
				Error: fmt.Sprintf("解析 %s 失败: %s", filepath.Base(file), err.Error()),
			})
			continue
		}

		presetName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		preset, err := frpc.ConvertToModels(config, presetName)
		if err != nil {
			results = append(results, ImportResult{
				Error: fmt.Sprintf("转换 %s 失败: %s", filepath.Base(file), err.Error()),
			})
			continue
		}

		results = append(results, ImportResult{Preset: preset})
	}

	return results
}

func (a *App) ExportPresetToml(presetName string, tomlContent string) string {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "导出 FrpEasy 预设",
		DefaultFilename: fmt.Sprintf("预设-%s.toml", presetName),
		Filters: []runtime.FileFilter{
			{DisplayName: "TOML 文件 (*.toml)", Pattern: "*.toml"},
		},
	})

	if err != nil {
		return ""
	}

	if file == "" {
		return ""
	}

	err = os.WriteFile(file, []byte(tomlContent), 0644)
	if err != nil {
		return ""
	}

	return file
}

func (a *App) ImportPresetFromToml() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "导入 FrpEasy 预设",
		Filters: []runtime.FileFilter{
			{DisplayName: "TOML 文件 (*.toml)", Pattern: "*.toml"},
		},
	})

	if err != nil {
		return ""
	}

	if file == "" {
		return ""
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	return string(content)
}

func (a *App) ExportPresetAsTomlBatch(serversJson string, servicesJson string, presetName string) string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 frp 配置导出目录",
	})

	if err != nil || dir == "" {
		return ""
	}

	var servers []struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Port    int    `json:"port"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal([]byte(serversJson), &servers); err != nil {
		return ""
	}

	var services []models.Service
	if err := json.Unmarshal([]byte(servicesJson), &services); err != nil {
		return ""
	}

	for _, s := range servers {
		server := models.Server{
			Name:    s.Name,
			Address: s.Address,
			Port:    s.Port,
			Token:   s.Token,
		}

		filename := fmt.Sprintf("frpc-%s.toml", sanitizeFilename(s.Name))
		filepath := filepath.Join(dir, filename)
		content := frpc.GenerateConfig(&server, services)
		os.WriteFile(filepath, []byte(content), 0644)
	}

	return dir
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	return re.ReplaceAllString(name, "_")
}

func (a *App) SaveAppConfig(tomlContent string) error {
	appConfig, err := config.ParseConfigString(tomlContent)
	if err != nil {
		fmt.Println("Failed to parse config TOML:", err)
		return fmt.Errorf("failed to parse config: %w", err)
	}
	fmt.Println("Saving config, presets count:", len(appConfig.Presets))
	return config.SaveConfig(filepath.Join(a.dataDir, "config.toml"), appConfig)
}

func (a *App) LoadAppConfig() string {
	appConfig, err := config.LoadConfig(filepath.Join(a.dataDir, "config.toml"))
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return ""
	}
	tomlStr, err := config.ToTomlString(appConfig)
	if err != nil {
		fmt.Println("Failed to marshal config:", err)
		return ""
	}
	fmt.Println("Loaded config, presets count:", len(appConfig.Presets))
	return tomlStr
}

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) GetLatestFrpcVersion() string {
	version, err := frpc.GetLatestFrpcVersion()
	if err != nil {
		fmt.Println("Failed to get latest frpc version:", err)
		return ""
	}
	return version
}

func (a *App) GetCurrentFrpcVersion() string {
	version, err := frpc.GetFrpcVersion(filepath.Join(a.dataDir, "bin"))
	if err != nil {
		return ""
	}
	return version
}

func (a *App) CompareFrpcVersions(v1, v2 string) int {
	return frpc.CompareVersions(v1, v2)
}
