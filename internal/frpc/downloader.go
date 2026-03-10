package frpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"frpeasy/internal/models"
)

const (
	FrpcVersion = "0.61.1"
)

type DownloadSource string

const (
	SourceGithub  DownloadSource = "github"
	SourceGHProxy DownloadSource = "ghproxy"
	SourceFastGit DownloadSource = "fastgit"
	SourceMoeyy   DownloadSource = "moeyy"
)

func GetDownloadURL(source DownloadSource) string {
	filename := GetFrpcFilename()
	ext := ".zip"
	if runtime.GOOS != "windows" {
		ext = ".tar.gz"
	}

	githubURL := fmt.Sprintf("https://github.com/fatedier/frp/releases/download/v%s/%s%s", FrpcVersion, filename, ext)

	switch source {
	case SourceGHProxy:
		return "https://ghproxy.net/" + githubURL
	case SourceFastGit:
		return "https://hub.fastgit.xyz/fatedier/frp/releases/download/v" + FrpcVersion + "/" + filename + ext
	case SourceMoeyy:
		return "https://github.moeyy.xyz/" + githubURL
	default:
		return githubURL
	}
}

func GetFrpcFilename() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var osName, archName string

	switch goos {
	case "windows":
		osName = "windows"
	case "darwin":
		osName = "darwin"
	case "linux":
		osName = "linux"
	default:
		osName = goos
	}

	switch goarch {
	case "amd64":
		archName = "amd64"
	case "arm64":
		archName = "arm64"
	case "386":
		archName = "386"
	case "arm":
		archName = "arm"
	default:
		archName = goarch
	}

	return fmt.Sprintf("frp_%s_%s_%s", FrpcVersion, osName, archName)
}

func GetFrpcExeName() string {
	if runtime.GOOS == "windows" {
		return "frpc.exe"
	}
	return "frpc"
}

func GetFrpcPath(binDir string) string {
	return filepath.Join(binDir, GetFrpcExeName())
}

func IsFrpcDownloaded(binDir string) bool {
	frpcPath := GetFrpcPath(binDir)
	_, err := os.Stat(frpcPath)
	return err == nil
}

func DownloadFrpc(binDir string, source DownloadSource, progress func(models.DownloadProgress)) error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	url := GetDownloadURL(source)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download frpc: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download frpc: HTTP %d", resp.StatusCode)
	}

	tempFile := filepath.Join(binDir, "frpc_download.tmp")
	if runtime.GOOS == "windows" {
		tempFile += ".zip"
	} else {
		tempFile += ".tar.gz"
	}

	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	total := resp.ContentLength
	downloaded := int64(0)
	buf := make([]byte, 32*1024)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				os.Remove(tempFile)
				return fmt.Errorf("failed to write file: %w", writeErr)
			}
			downloaded += int64(n)
			if progress != nil && total > 0 {
				progress(models.DownloadProgress{
					TotalBytes: total,
					Downloaded: downloaded,
					Percentage: float64(downloaded) / float64(total) * 100,
					IsComplete: false,
				})
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(tempFile)
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	out.Close()

	if err := extractFrpc(binDir, tempFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to extract frpc: %w", err)
	}

	os.Remove(tempFile)

	if progress != nil {
		progress(models.DownloadProgress{
			TotalBytes: total,
			Downloaded: total,
			Percentage: 100,
			IsComplete: true,
		})
	}

	return nil
}

func extractFrpc(binDir, archivePath string) error {
	frpcExe := GetFrpcExeName()
	targetPath := GetFrpcPath(binDir)

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", archivePath, binDir))
		hideWindow(cmd)
		if err := cmd.Run(); err != nil {
			return err
		}

		matches, _ := filepath.Glob(filepath.Join(binDir, "frp_*", frpcExe))
		if len(matches) > 0 {
			os.Rename(matches[0], targetPath)
		}
	} else {
		cmd := exec.Command("tar", "-xzf", archivePath, "-C", binDir)
		if err := cmd.Run(); err != nil {
			return err
		}

		matches, _ := filepath.Glob(filepath.Join(binDir, "frp_*", frpcExe))
		if len(matches) > 0 {
			os.Rename(matches[0], targetPath)
		}

		os.Chmod(targetPath, 0755)
	}

	cleanExtractDir(binDir)

	return nil
}

func cleanExtractDir(binDir string) {
	entries, _ := os.ReadDir(binDir)
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "frp_") {
			os.RemoveAll(filepath.Join(binDir, entry.Name()))
		}
	}
}

func GetFrpcVersion(binDir string) (string, error) {
	frpcPath := GetFrpcPath(binDir)
	if _, err := os.Stat(frpcPath); os.IsNotExist(err) {
		return "", fmt.Errorf("frpc not found")
	}

	cmd := exec.Command(frpcPath, "-v")
	hideWindow(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func GetLatestFrpcVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/fatedier/frp/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest version: HTTP %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

func CompareVersions(v1, v2 string) int {
	re := regexp.MustCompile(`[^0-9.]`)
	v1 = re.ReplaceAllString(v1, "")
	v2 = re.ReplaceAllString(v2, "")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}

	return 0
}
