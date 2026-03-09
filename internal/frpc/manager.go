package frpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"syscall"
	"time"

	"frpeasy/internal/models"

	"github.com/google/uuid"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func hideWindow(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}
}

type ProcessManager struct {
	mu        sync.RWMutex
	processes map[string]*ProcessInfo
	binDir    string
	configDir string
}

type ProcessInfo struct {
	PresetID   string
	ServerID   string
	Cmd        *exec.Cmd
	ConfigPath string
	Logs       []models.LogEntry
	StartTime  time.Time
	Cancel     context.CancelFunc
}

func NewProcessManager(binDir, configDir string) *ProcessManager {
	return &ProcessManager{
		processes: make(map[string]*ProcessInfo),
		binDir:    binDir,
		configDir: configDir,
	}
}

func (pm *ProcessManager) getKey(presetID, serverID string) string {
	return fmt.Sprintf("%s:%s", presetID, serverID)
}

func (pm *ProcessManager) IsRunning(presetID, serverID string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	key := pm.getKey(presetID, serverID)
	info, exists := pm.processes[key]
	if !exists {
		return false
	}

	if info.Cmd == nil || info.Cmd.Process == nil {
		return false
	}

	return info.Cmd.ProcessState == nil
}

func (pm *ProcessManager) Start(presetID, serverID string, server *models.Server, services []models.Service, onLog func(presetID, serverID string, log models.LogEntry)) error {
	if pm.IsRunning(presetID, serverID) {
		return fmt.Errorf("server is already running")
	}

	if !IsFrpcDownloaded(pm.binDir) {
		return fmt.Errorf("frpc binary not found, please download first")
	}

	configPath, err := WriteConfigFile(pm.configDir, presetID, serverID, server, services)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	frpcPath := GetFrpcPath(pm.binDir)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, frpcPath, "-c", configPath)
	hideWindow(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		DeleteConfigFile(configPath)
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		DeleteConfigFile(configPath)
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		DeleteConfigFile(configPath)
		return fmt.Errorf("failed to start frpc: %w", err)
	}

	info := &ProcessInfo{
		PresetID:   presetID,
		ServerID:   serverID,
		Cmd:        cmd,
		ConfigPath: configPath,
		Logs:       make([]models.LogEntry, 0),
		StartTime:  time.Now(),
		Cancel:     cancel,
	}

	pm.mu.Lock()
	pm.processes[pm.getKey(presetID, serverID)] = info
	pm.mu.Unlock()

	go pm.readLogs(info, stdout, stderr, onLog)
	go pm.waitProcess(presetID, serverID)

	return nil
}

func (pm *ProcessManager) readLogs(info *ProcessInfo, stdout, stderr io.Reader, onLog func(string, string, models.LogEntry)) {
	generateID := func() string {
		return uuid.New().String()[:8]
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := stripANSI(scanner.Text())
		logEntry := models.LogEntry{
			ID:        generateID(),
			Timestamp: time.Now().UnixMilli(),
			Message:   line,
			Type:      "info",
		}

		pm.mu.Lock()
		info.Logs = append(info.Logs, logEntry)
		if len(info.Logs) > 100 {
			info.Logs = info.Logs[1:]
		}
		pm.mu.Unlock()

		if onLog != nil {
			onLog(info.PresetID, info.ServerID, logEntry)
		}
	}

	scannerErr := bufio.NewScanner(stderr)
	for scannerErr.Scan() {
		line := stripANSI(scannerErr.Text())
		logEntry := models.LogEntry{
			ID:        generateID(),
			Timestamp: time.Now().UnixMilli(),
			Message:   line,
			Type:      "error",
		}

		pm.mu.Lock()
		info.Logs = append(info.Logs, logEntry)
		if len(info.Logs) > 100 {
			info.Logs = info.Logs[1:]
		}
		pm.mu.Unlock()

		if onLog != nil {
			onLog(info.PresetID, info.ServerID, logEntry)
		}
	}
}

func (pm *ProcessManager) waitProcess(presetID, serverID string) {
	key := pm.getKey(presetID, serverID)

	pm.mu.RLock()
	info, exists := pm.processes[key]
	pm.mu.RUnlock()

	if !exists || info.Cmd == nil {
		return
	}

	info.Cmd.Wait()
}

func (pm *ProcessManager) Stop(presetID, serverID string) error {
	key := pm.getKey(presetID, serverID)

	pm.mu.Lock()
	defer pm.mu.Unlock()

	info, exists := pm.processes[key]
	if !exists {
		return nil
	}

	if info.Cancel != nil {
		info.Cancel()
	}

	if info.Cmd != nil && info.Cmd.Process != nil {
		if info.Cmd.ProcessState == nil {
			if runtime.GOOS == "windows" {
				cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", info.Cmd.Process.Pid))
				hideWindow(cmd)
				cmd.Run()
			} else {
				info.Cmd.Process.Signal(syscall.SIGTERM)
				time.Sleep(1 * time.Second)
				if info.Cmd.ProcessState == nil {
					info.Cmd.Process.Kill()
				}
			}
		}
	}

	if info.ConfigPath != "" {
		DeleteConfigFile(info.ConfigPath)
	}

	delete(pm.processes, key)
	return nil
}

func (pm *ProcessManager) StopAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for key, info := range pm.processes {
		if info.Cancel != nil {
			info.Cancel()
		}
		if info.Cmd != nil && info.Cmd.Process != nil && info.Cmd.ProcessState == nil {
			if runtime.GOOS == "windows" {
				cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", info.Cmd.Process.Pid))
				hideWindow(cmd)
				cmd.Run()
			} else {
				info.Cmd.Process.Signal(syscall.SIGTERM)
			}
		}
		if info.ConfigPath != "" {
			DeleteConfigFile(info.ConfigPath)
		}
		delete(pm.processes, key)
	}
}

func (pm *ProcessManager) GetLogs(presetID, serverID string) []models.LogEntry {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	key := pm.getKey(presetID, serverID)
	if info, exists := pm.processes[key]; exists {
		logs := make([]models.LogEntry, len(info.Logs))
		copy(logs, info.Logs)
		return logs
	}
	return nil
}

func (pm *ProcessManager) GetUptime(presetID, serverID string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	key := pm.getKey(presetID, serverID)
	if info, exists := pm.processes[key]; exists {
		return int(time.Since(info.StartTime).Seconds())
	}
	return 0
}

func (pm *ProcessManager) GetRunningServers() []models.ServerRuntime {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make([]models.ServerRuntime, 0, len(pm.processes))
	for _, info := range pm.processes {
		pid := 0
		if info.Cmd != nil && info.Cmd.Process != nil {
			pid = info.Cmd.Process.Pid
		}
		result = append(result, models.ServerRuntime{
			PresetID:   info.PresetID,
			ServerID:   info.ServerID,
			ProcessPID: pid,
			ConfigPath: info.ConfigPath,
		})
	}
	return result
}
