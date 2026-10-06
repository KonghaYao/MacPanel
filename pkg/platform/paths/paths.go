package paths

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const macPanelDirName = "MacPanel"

func BaseDir() string {
	if runtime.GOOS == "darwin" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("resolve home dir: %v", err))
		}
		return filepath.Join(home, "Library", "Application Support", macPanelDirName)
	}
	return ""
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
func Bootstrap(version string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

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
		version = "v2.0.0"
	}
	password, err := generatePassword(16)
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`BASE_DIR=%s
ORIGINAL_PORT=9999
ORIGINAL_VERSION=%s
ORIGINAL_USERNAME=admin
ORIGINAL_PASSWORD=%s
ORIGINAL_ENTRANCE=
LANGUAGE=zh
PANEL_EDITION=standard
`, DataDir(), version, password)
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
