package service

import "github.com/1Panel-dev/1Panel/pkg/platform/capabilities"

func rejectDarwin() error {
	return capabilities.RejectDarwin()
}

func rejectDarwinFeature(feature string) error {
	return capabilities.RejectDarwinUnless(feature)
}

func rejectDarwinDockerServiceControl() error {
	if capabilities.IsDarwin() {
		return capabilities.ErrNotSupportedOnMac
	}
	return nil
}
