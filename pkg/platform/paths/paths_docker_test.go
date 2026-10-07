package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveDockerDaemonJsonPathSnap(t *testing.T) {
	got := ResolveDockerDaemonJsonPath("/snap/bin/docker")
	if got != snapDockerDaemonJSONPath {
		t.Fatalf("snap path = %q, want %q", got, snapDockerDaemonJSONPath)
	}
}

func TestResolveDockerDaemonJsonPathDefault(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("darwin uses home directory path")
	}
	got := ResolveDockerDaemonJsonPath("/usr/bin/docker")
	if got != defaultDockerDaemonJSONPath {
		t.Fatalf("default path = %q, want %q", got, defaultDockerDaemonJSONPath)
	}
}

func TestResolveDarwinDockerDaemonJSONPathDockerDesktop(t *testing.T) {
	home := t.TempDir()
	got := resolveDarwinDockerDaemonJSONPath("/usr/local/bin/docker", home)
	want := filepath.Join(home, ".docker", "daemon.json")
	if got != want {
		t.Fatalf("docker desktop path = %q, want %q", got, want)
	}
}

func TestResolveDarwinDockerDaemonJSONPathOrbStackBinary(t *testing.T) {
	home := t.TempDir()
	dockerPath := filepath.Join(home, ".orbstack", "bin", "docker")
	got := resolveDarwinDockerDaemonJSONPath(dockerPath, home)
	want := filepath.Join(home, ".orbstack", "config", "docker.json")
	if got != want {
		t.Fatalf("orbstack binary path = %q, want %q", got, want)
	}
}

func TestResolveDarwinDockerDaemonJSONPathOrbStackContext(t *testing.T) {
	home := t.TempDir()
	dockerConfigDir := filepath.Join(home, ".docker")
	if err := os.MkdirAll(dockerConfigDir, 0o755); err != nil {
		t.Fatalf("mkdir .docker: %v", err)
	}
	config := []byte(`{"currentContext":"orbstack"}`)
	if err := os.WriteFile(filepath.Join(dockerConfigDir, "config.json"), config, 0o644); err != nil {
		t.Fatalf("write docker config: %v", err)
	}

	got := resolveDarwinDockerDaemonJSONPath("/usr/local/bin/docker", home)
	want := filepath.Join(home, ".orbstack", "config", "docker.json")
	if got != want {
		t.Fatalf("orbstack context path = %q, want %q", got, want)
	}
}

func TestResolveDarwinDockerDaemonJSONPathOrbStackConfigExists(t *testing.T) {
	home := t.TempDir()
	orbConfigDir := filepath.Join(home, ".orbstack", "config")
	orbBinDir := filepath.Join(home, ".orbstack", "bin")
	if err := os.MkdirAll(orbConfigDir, 0o755); err != nil {
		t.Fatalf("mkdir orbstack config: %v", err)
	}
	if err := os.MkdirAll(orbBinDir, 0o755); err != nil {
		t.Fatalf("mkdir orbstack bin: %v", err)
	}
	if err := os.WriteFile(filepath.Join(orbConfigDir, "docker.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("write orbstack docker.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(orbBinDir, "docker"), []byte(""), 0o755); err != nil {
		t.Fatalf("write orbstack docker binary: %v", err)
	}

	got := resolveDarwinDockerDaemonJSONPath("/usr/local/bin/docker", home)
	want := filepath.Join(home, ".orbstack", "config", "docker.json")
	if got != want {
		t.Fatalf("orbstack config path = %q, want %q", got, want)
	}
}

func TestDetectDarwinDockerRuntimeOrbStackEnvContext(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DOCKER_CONTEXT", "orbstack")
	got := detectDarwinDockerRuntime("/usr/local/bin/docker", home)
	if got != DockerRuntimeOrbStack {
		t.Fatalf("runtime = %q, want %q", got, DockerRuntimeOrbStack)
	}
}

func TestDetectDarwinDockerRuntimeDockerDesktop(t *testing.T) {
	home := t.TempDir()
	got := detectDarwinDockerRuntime("/usr/local/bin/docker", home)
	if got != DockerRuntimeDockerDesktop {
		t.Fatalf("runtime = %q, want %q", got, DockerRuntimeDockerDesktop)
	}
}

func TestResolveDockerDaemonJsonPathDarwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only test")
	}
	got := ResolveDockerDaemonJsonPath("/usr/local/bin/docker")
	if filepath.Base(filepath.Dir(got)) != ".docker" || filepath.Base(got) != "daemon.json" {
		t.Fatalf("darwin path = %q, want ~/.docker/daemon.json shape", got)
	}
}
