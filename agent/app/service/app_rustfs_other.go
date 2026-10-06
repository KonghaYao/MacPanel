//go:build !darwin

package service

import "github.com/1Panel-dev/1Panel/agent/app/model"

func patchRustFSComposeForInstallPlatform(appInstall *model.AppInstall) error {
	return nil
}

func patchRustFSInitScriptForInstallPlatform(appInstall *model.AppInstall) error {
	return nil
}
