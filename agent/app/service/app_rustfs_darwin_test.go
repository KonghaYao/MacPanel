//go:build darwin

package service

import (
	"strings"
	"testing"
)

const rustfsSampleCompose = `services:
  rustfs:
    image: rustfs/rustfs:1.0.0
    container_name: ${CONTAINER_NAME}
    volumes:
      - ./data:/data
networks:
  1panel-network:
    external: true
`

func TestPatchRustFSComposeYAML(t *testing.T) {
	patched, err := patchRustFSComposeYAML(rustfsSampleCompose)
	if err != nil {
		t.Fatalf("patch compose: %v", err)
	}
	if strings.Contains(patched, "./data:/data") {
		t.Fatalf("expected bind mount to be replaced, got %q", patched)
	}
	if !strings.Contains(patched, rustfsNamedVolumeRef) {
		t.Fatalf("expected named volume mount %q in %q", rustfsNamedVolumeRef, patched)
	}
	if !strings.Contains(patched, rustfsNamedVolumeName+":") {
		t.Fatalf("expected top-level named volume %q in %q", rustfsNamedVolumeName, patched)
	}
}

func TestIsRustFSDataBindMount(t *testing.T) {
	cases := map[string]bool{
		"./data:/data":     true,
		"./data:/data:rw":  true,
		"data:/data":       true,
		"./logs:/logs":     false,
		"rustfs-data:/data": false,
	}
	for volume, want := range cases {
		if got := isRustFSDataBindMount(volume); got != want {
			t.Fatalf("isRustFSDataBindMount(%q) = %v, want %v", volume, got, want)
		}
	}
}
