package mirrors

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func readPip(opts Options) (map[string]string, error) {
	path, err := pipConfigPath(opts)
	if err != nil {
		return nil, err
	}
	current := map[string]string{}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return current, nil
		}
		return nil, err
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	section := cfg.Section("global")
	if section.HasKey("index-url") {
		current["indexUrl"] = section.Key("index-url").String()
	}
	if section.HasKey("trusted-host") {
		current["trustedHost"] = section.Key("trusted-host").String()
	}
	return current, nil
}

func writePip(opts Options, values map[string]string) error {
	path, err := pipConfigPath(opts)
	if err != nil {
		return err
	}
	cfg := ini.Empty()
	if _, err := os.Stat(path); err == nil {
		loaded, loadErr := ini.Load(path)
		if loadErr != nil {
			return loadErr
		}
		cfg = loaded
	} else if !os.IsNotExist(err) {
		return err
	}
	section := cfg.Section("global")
	setOrDeleteIni(section, "index-url", values["indexUrl"])
	setOrDeleteIni(section, "trusted-host", values["trustedHost"])
	if !iniHasKeys(cfg) {
		return removeFile(path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".macpanel-mirror-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpName)
	if err := cfg.SaveTo(tmpName); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func setOrDeleteIni(section *ini.Section, key, value string) {
	if value == "" {
		section.DeleteKey(key)
		return
	}
	section.Key(key).SetValue(value)
}

func iniHasKeys(cfg *ini.File) bool {
	for _, section := range cfg.Sections() {
		if len(section.Keys()) > 0 {
			return true
		}
	}
	return false
}
