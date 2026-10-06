package service

import "errors"

var ErrDesktopWallpaperUnsupported = errors.New("desktop wallpaper is unsupported on this platform")

type DesktopService struct{}

type IDesktopService interface {
	GetCompressedWallpaper() ([]byte, error)
}

func NewIDesktopService() IDesktopService {
	return &DesktopService{}
}
