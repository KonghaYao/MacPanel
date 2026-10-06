package capabilities

import (
	"runtime"
)

type FeatureMap struct {
	Firewall          bool `json:"firewall"`
	Fail2ban          bool `json:"fail2ban"`
	DiskManagement    bool `json:"disk_management"`
	Fstab             bool `json:"fstab"`
	Swap              bool `json:"swap"`
	NtpSync           bool `json:"ntp_sync"`
	SshdConfig        bool `json:"sshd_config"`
	DockerInstall     bool `json:"docker_install"`
	DockerManage      bool `json:"docker_manage"`
	OnlineUpgrade     bool `json:"online_upgrade"`
	ProcMonitoring    bool `json:"proc_monitoring"`
	OpenrestyDiagnose bool `json:"openresty_diagnose"`
}

type PlatformCapabilities struct {
	OS       string     `json:"os"`
	Arch     string     `json:"arch"`
	Features FeatureMap `json:"features"`
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func RejectDarwin() error {
	if IsDarwin() {
		return ErrNotSupportedOnMac
	}
	return nil
}

func RejectDarwinUnless(feature string) error {
	if !IsDarwin() {
		return nil
	}
	if isFeatureEnabled(feature) {
		return nil
	}
	return ErrNotSupportedOnMac
}

func Current() PlatformCapabilities {
	if IsDarwin() {
		return PlatformCapabilities{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Features: darwinFeatures(),
		}
	}
	return PlatformCapabilities{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Features: linuxFeatures(),
	}
}

func linuxFeatures() FeatureMap {
	return FeatureMap{
		Firewall:          true,
		Fail2ban:          true,
		DiskManagement:    true,
		Fstab:             true,
		Swap:              true,
		NtpSync:           true,
		SshdConfig:        true,
		DockerInstall:     true,
		DockerManage:      true,
		OnlineUpgrade:     true,
		ProcMonitoring:    true,
		OpenrestyDiagnose: true,
	}
}

func darwinFeatures() FeatureMap {
	return FeatureMap{
		Firewall:          false,
		Fail2ban:          false,
		DiskManagement:    false,
		Fstab:             false,
		Swap:              false,
		NtpSync:           false,
		SshdConfig:        false,
		DockerInstall:     false,
		DockerManage:      true,
		OnlineUpgrade:     false,
		ProcMonitoring:    true,
		OpenrestyDiagnose: false,
	}
}

func isFeatureEnabled(feature string) bool {
	features := darwinFeatures()
	switch feature {
	case "firewall":
		return features.Firewall
	case "fail2ban":
		return features.Fail2ban
	case "disk_management":
		return features.DiskManagement
	case "fstab":
		return features.Fstab
	case "swap":
		return features.Swap
	case "ntp_sync":
		return features.NtpSync
	case "sshd_config":
		return features.SshdConfig
	case "docker_install":
		return features.DockerInstall
	case "docker_manage":
		return features.DockerManage
	case "online_upgrade":
		return features.OnlineUpgrade
	case "proc_monitoring":
		return features.ProcMonitoring
	case "openresty_diagnose":
		return features.OpenrestyDiagnose
	default:
		return false
	}
}
