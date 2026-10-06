package service

import "github.com/1Panel-dev/1Panel/agent/app/model"

const rustfsAppKey = "rustfs"

func patchRustFSComposeForInstall(app model.App, appInstall *model.AppInstall) error {
	if app.Key != rustfsAppKey {
		return nil
	}
	return patchRustFSComposeForInstallPlatform(appInstall)
}

func patchRustFSInitScriptForInstall(app model.App, appInstall *model.AppInstall) error {
	if app.Key != rustfsAppKey {
		return nil
	}
	return patchRustFSInitScriptForInstallPlatform(appInstall)
}
