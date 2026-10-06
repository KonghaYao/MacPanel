//go:build darwin

package service

import (
	"path"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/utils/files"
	"gopkg.in/yaml.v3"
)

const (
	rustfsDataMountTarget  = "/data"
	rustfsNamedVolumeSuffix = "-data"
	rustfsDarwinInitScript = "#!/bin/bash\n\n# RustFS data uses a Docker named volume on macOS; host chown is unsupported.\n"
)

func rustfsNamedVolumeName(containerName string) string {
	return containerName + rustfsNamedVolumeSuffix
}

func patchRustFSComposeForInstallPlatform(appInstall *model.AppInstall) error {
	patched, err := patchRustFSComposeYAML(appInstall.DockerCompose, appInstall.ContainerName)
	if err != nil {
		return err
	}
	appInstall.DockerCompose = patched
	return nil
}

func patchRustFSInitScriptForInstallPlatform(appInstall *model.AppInstall) error {
	initPath := path.Join(appInstall.GetPath(), "scripts", "init.sh")
	fileOp := files.NewFileOp()
	if !fileOp.Stat(initPath) {
		return nil
	}
	return fileOp.WriteFile(initPath, strings.NewReader(rustfsDarwinInitScript), constant.DirPerm)
}

func patchRustFSComposeYAML(compose string, containerName string) (string, error) {
	if containerName == "" {
		return compose, nil
	}
	volumeName := rustfsNamedVolumeName(containerName)
	volumeRef := volumeName + ":" + rustfsDataMountTarget

	var composeMap map[string]interface{}
	if err := yaml.Unmarshal([]byte(compose), &composeMap); err != nil {
		return "", err
	}
	services, ok := composeMap["services"].(map[string]interface{})
	if !ok {
		return compose, nil
	}

	patched := false
	for _, service := range services {
		serviceMap, ok := service.(map[string]interface{})
		if !ok {
			continue
		}
		volumes, ok := serviceMap["volumes"].([]interface{})
		if !ok {
			continue
		}
		for i, volume := range volumes {
			volumeStr, ok := volume.(string)
			if !ok || !isRustFSDataBindMount(volumeStr) {
				continue
			}
			volumes[i] = volumeRef
			patched = true
		}
		serviceMap["volumes"] = volumes
	}
	if !patched {
		return compose, nil
	}

	topVolumes, _ := composeMap["volumes"].(map[string]interface{})
	if topVolumes == nil {
		topVolumes = make(map[string]interface{})
	}
	topVolumes[volumeName] = map[string]interface{}{}
	composeMap["volumes"] = topVolumes

	out, err := yaml.Marshal(composeMap)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isRustFSDataBindMount(volume string) bool {
	parts := strings.SplitN(volume, ":", 3)
	if len(parts) < 2 {
		return false
	}
	source := strings.TrimSuffix(strings.TrimSpace(parts[0]), "/")
	target := strings.TrimSpace(parts[1])
	return (source == "./data" || source == "data") && strings.HasPrefix(target, rustfsDataMountTarget)
}
