//go:build !darwin

package service

func (s *DesktopService) GetCompressedWallpaper() ([]byte, error) {
	return nil, ErrDesktopWallpaperUnsupported
}
