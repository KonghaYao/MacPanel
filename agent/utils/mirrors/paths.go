package mirrors

import (
	"path/filepath"
)

type Options struct {
	Home             string
	GOOS             string
	GoEnvFile        string
	XDGConfigHome    string
	CargoHome        string
	PipConfigFile    string
	NpmrcFile        string
	GradleUserHome   string
	DockerDaemonPath string
}

func (o Options) validate() error {
	if o.Home == "" {
		return invalidf("home directory is empty")
	}
	if o.GOOS == "" {
		return invalidf("GOOS is empty")
	}
	if o.DockerDaemonPath == "" {
		return invalidf("docker daemon.json path is empty")
	}
	return nil
}

func userConfigRoot(opts Options) string {
	if opts.GOOS == "darwin" {
		return filepath.Join(opts.Home, "Library", "Application Support")
	}
	if opts.XDGConfigHome != "" {
		return opts.XDGConfigHome
	}
	return filepath.Join(opts.Home, ".config")
}

func pipConfigPath(opts Options) (string, error) {
	if opts.PipConfigFile != "" {
		return opts.PipConfigFile, nil
	}
	return filepath.Join(userConfigRoot(opts), "pip", "pip.conf"), nil
}

func npmrcPath(opts Options) (string, error) {
	if opts.NpmrcFile != "" {
		return opts.NpmrcFile, nil
	}
	return filepath.Join(opts.Home, ".npmrc"), nil
}

func mavenSettingsPath(opts Options) (string, error) {
	return filepath.Join(opts.Home, ".m2", "settings.xml"), nil
}

func gradleInitPath(opts Options) (string, error) {
	root := opts.GradleUserHome
	if root == "" {
		root = filepath.Join(opts.Home, ".gradle")
	}
	return filepath.Join(root, "init.d", "macpanel-mirror.init.gradle"), nil
}

func goEnvPath(opts Options) (string, error) {
	if opts.GoEnvFile != "" {
		return opts.GoEnvFile, nil
	}
	return filepath.Join(userConfigRoot(opts), "go", "env"), nil
}

func cargoConfigPath(opts Options) (string, error) {
	root := opts.CargoHome
	if root == "" {
		root = filepath.Join(opts.Home, ".cargo")
	}
	return filepath.Join(root, "config.toml"), nil
}

func dockerDaemonPath(opts Options) (string, error) {
	return opts.DockerDaemonPath, nil
}
