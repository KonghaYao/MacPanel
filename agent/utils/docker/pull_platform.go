package docker

import (
	"strings"

	"github.com/docker/docker/api/types/image"
)

// ApplyPullPlatform sets moby ImagePull platform when a concrete value is requested.
// Empty or "auto" leaves the daemon default (host architecture).
func ApplyPullPlatform(options *image.PullOptions, platform string) {
	platform = strings.TrimSpace(platform)
	if platform == "" || platform == "auto" {
		return
	}
	options.Platform = platform
}

func NewPullOptions(imageName, platform string) image.PullOptions {
	options := image.PullOptions{}
	if authStr, ok := loadRegistryAuthFromDockerConfig(imageName); ok {
		options.RegistryAuth = authStr
	}
	ApplyPullPlatform(&options, platform)
	return options
}
