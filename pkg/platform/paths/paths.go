package paths

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	panelversion "github.com/1Panel-dev/1Panel/pkg/platform/version"
	"runtime"
	"strings"
	"sync"
)

var bootstrapOnce sync.Once
var bootstrapErr error

// ResetBootstrapState clears one-time bootstrap state. For tests only.
func ResetBootstrapState() {
	bootstrapOnce = sync.Once{}
	bootstrapErr = nil
}

const macPanelDirName = "MacPanel"

const (
	defaultDockerDaemonJSONPath = "/etc/docker/daemon.json"
	snapDockerDaemonJSONPath    = "/var/snap/docker/current/config/daemon.json"
)

type DockerRuntime string

const (
	DockerRuntimeLinux         DockerRuntime = "linux"
	DockerRuntimeDockerDesktop DockerRuntime = "docker-desktop"
	DockerRuntimeOrbStack      DockerRuntime = "orbstack"
)

func ResolveDockerDaemonJsonPath(dockerBinaryPath string) string {
	if strings.Contains(dockerBinaryPath, "snap") {
		return snapDockerDaemonJSONPath
	}
	if runtime.GOOS == "darwin" {
		if home, err := os.UserHomeDir(); err == nil {
			return resolveDarwinDockerDaemonJSONPath(dockerBinaryPath, home)
		}
	}
	return defaultDockerDaemonJSONPath
}

func DetectDockerRuntime(dockerBinaryPath string) DockerRuntime {
	if runtime.GOOS != "darwin" {
		return DockerRuntimeLinux
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return DockerRuntimeDockerDesktop
	}
	return detectDarwinDockerRuntime(dockerBinaryPath, home)
}

func resolveDarwinDockerDaemonJSONPath(dockerBinaryPath, home string) string {
	if detectDarwinDockerRuntime(dockerBinaryPath, home) == DockerRuntimeOrbStack {
		return filepath.Join(home, ".orbstack", "config", "docker.json")
	}
	return filepath.Join(home, ".docker", "daemon.json")
}

func detectDarwinDockerRuntime(dockerBinaryPath, home string) DockerRuntime {
	if isOrbStackDockerBinary(dockerBinaryPath) {
		return DockerRuntimeOrbStack
	}
	if readDockerCurrentContext(home) == "orbstack" {
		return DockerRuntimeOrbStack
	}
	orbConfig := filepath.Join(home, ".orbstack", "config", "docker.json")
	orbDocker := filepath.Join(home, ".orbstack", "bin", "docker")
	if fileExists(orbConfig) && fileExists(orbDocker) && !isDockerDesktopBinary(dockerBinaryPath) {
		return DockerRuntimeOrbStack
	}
	return DockerRuntimeDockerDesktop
}

func isOrbStackDockerBinary(dockerBinaryPath string) bool {
	normalized := filepath.ToSlash(dockerBinaryPath)
	return strings.Contains(normalized, "/.orbstack/") || strings.Contains(normalized, "orbstack")
}

func isDockerDesktopBinary(dockerBinaryPath string) bool {
	normalized := filepath.ToSlash(dockerBinaryPath)
	return strings.Contains(normalized, "Docker.app")
}

func readDockerCurrentContext(home string) string {
	if ctx := strings.TrimSpace(os.Getenv("DOCKER_CONTEXT")); ctx != "" {
		return ctx
	}
	configPath := filepath.Join(home, ".docker", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	var cfg struct {
		CurrentContext string `json:"currentContext"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	return strings.TrimSpace(cfg.CurrentContext)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func BaseDir() string {
	if runtime.GOOS == "darwin" {
		if override := os.Getenv("MACPANEL_HOME"); override != "" {
			return override
		}
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("resolve home dir: %v", err))
		}
		return filepath.Join(home, "Library", "Application Support", macPanelDirName)
	}
	return ""
}

func defaultPort() string {
	if port := os.Getenv("MACPANEL_PORT"); port != "" {
		return port
	}
	return "9999"
}

func ConfigFile() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "config", "1pctl")
	}
	return "/usr/local/bin/1pctl"
}

func BinDir() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "bin")
	}
	return "/usr/local/bin"
}

func CoreBinaryPath() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BinDir(), "macpanel")
	}
	return filepath.Join(BinDir(), "1panel-core")
}

func AgentBinaryPath() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BinDir(), "macpanel")
	}
	return filepath.Join(BinDir(), "1panel-agent")
}

func LangDir() string {
	return filepath.Join(BinDir(), "lang")
}

func SocketDir() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "run")
	}
	return "/etc/1panel"
}

func SocketPath() string {
	return filepath.Join(SocketDir(), "agent.sock")
}

func NodeProxyIDPath() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "run", ".nodeProxyID")
	}
	return "/etc/1panel/.nodeProxyID"
}

func DataDir() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "1panel")
	}
	return LinuxDefaultBaseDir()
}

func ConfDir() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "1panel", "conf")
	}
	return filepath.Join(LinuxDefaultBaseDir(), "conf")
}

func LinuxDefaultBaseDir() string {
	return "/opt/1panel"
}

func LinuxDevConfFile() string {
	return filepath.Join(LinuxDefaultBaseDir(), "conf", "app.yaml")
}

func RunDir() string {
	return SocketDir()
}

func PidFile() string {
	return filepath.Join(RunDir(), "macpanel.pid")
}

func DaemonLogFile() string {
	return filepath.Join(RunDir(), "macpanel.log")
}

func ConfigDir() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(BaseDir(), "config")
	}
	return BinDir()
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func RecycleBinDir(baseDir string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(baseDir, "1panel", ".1panel_clash")
	}
	return "/.1panel_clash"
}

// Bootstrap ensures macOS base directories and a default 1pctl config exist.
// It is idempotent and safe to call from any startup path (CLI init, main, viper).
func Bootstrap(version string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	bootstrapOnce.Do(func() {
		bootstrapErr = bootstrap(version)
	})
	return bootstrapErr
}

func bootstrap(version string) error {
	dirs := []string{
		ConfigDir(),
		RunDir(),
		ConfDir(),
		filepath.Join(DataDir(), "db"),
		filepath.Join(DataDir(), "log"),
		filepath.Join(DataDir(), "apps"),
		filepath.Join(DataDir(), "runtime"),
		filepath.Join(DataDir(), "backup"),
		filepath.Join(DataDir(), "cache"),
		filepath.Join(DataDir(), "geo"),
		BinDir(),
		LangDir(),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create dir %s: %w", dir, err)
		}
	}

	if _, err := os.Stat(ConfigFile()); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config %s: %w", ConfigFile(), err)
	}

	if version == "" {
		version = panelversion.Version
	}
	password, err := generatePassword(16)
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`BASE_DIR=%s
ORIGINAL_PORT=%s
ORIGINAL_VERSION=%s
ORIGINAL_USERNAME=admin
ORIGINAL_PASSWORD=%s
ORIGINAL_ENTRANCE=
LANGUAGE=zh
PANEL_EDITION=standard
`, DataDir(), defaultPort(), version, password)
	if err := os.WriteFile(ConfigFile(), []byte(content), 0o600); err != nil {
		return fmt.Errorf("write config %s: %w", ConfigFile(), err)
	}
	return nil
}

func generatePassword(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
