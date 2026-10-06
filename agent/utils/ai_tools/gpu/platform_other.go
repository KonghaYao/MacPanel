//go:build !darwin

package gpu

func platformProviders() []provider {
	return nil
}
