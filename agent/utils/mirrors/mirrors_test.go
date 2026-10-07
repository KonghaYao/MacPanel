package mirrors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testOpts(t *testing.T) Options {
	t.Helper()
	home := t.TempDir()
	return Options{
		Home:             home,
		GOOS:             "linux",
		DockerDaemonPath: filepath.Join(home, "daemon.json"),
	}
}

func TestCatalogPresets(t *testing.T) {
	defs := append(definitions(), homebrewDefinition())
	for _, def := range defs {
		seen := map[string]bool{}
		hasOfficial := false
		for _, preset := range def.presets {
			if seen[preset.id] {
				t.Fatalf("%s duplicate preset %s", def.id, preset.id)
			}
			seen[preset.id] = true
			if preset.id == "official" {
				hasOfficial = true
			}
			values := completeValues(def, preset.values)
			if err := validateValues(def, values); err != nil {
				t.Fatalf("%s preset %s: %v", def.id, preset.id, err)
			}
		}
		if !hasOfficial {
			t.Fatalf("%s missing official preset", def.id)
		}
	}
}

func TestApplySwitchAndPreserve(t *testing.T) {
	opts := testOpts(t)

	if err := Apply(opts, ApplyRequest{Ecosystem: "pip", PresetID: "tsinghua"}); err != nil {
		t.Fatal(err)
	}
	pipPath, _ := pipConfigPath(opts)
	extra := "\n[install]\nuser = true\n"
	original, err := os.ReadFile(pipPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pipPath, append(original, []byte(extra)...), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "pip", PresetID: "aliyun"}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(pipPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "https://mirrors.aliyun.com/pypi/simple/") || !strings.Contains(text, "user = true") {
		t.Fatalf("pip config = %s", text)
	}

	npmPath, _ := npmrcPath(opts)
	if err := os.WriteFile(npmPath, []byte("//registry.npmjs.org/:_authToken=abc\nregistry=https://registry.npmjs.org/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "npm", PresetID: "npmmirror"}); err != nil {
		t.Fatal(err)
	}
	npmBody, _ := os.ReadFile(npmPath)
	if !strings.Contains(string(npmBody), "_authToken=abc") || !strings.Contains(string(npmBody), "registry=https://registry.npmmirror.com") {
		t.Fatalf("npmrc = %s", npmBody)
	}

	settings := filepath.Join(opts.Home, ".m2", "settings.xml")
	if err := os.MkdirAll(filepath.Dir(settings), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settings, []byte("<settings>\n  <servers><server><id>private</id></server></servers>\n</settings>\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "maven", PresetID: "aliyun"}); err != nil {
		t.Fatal(err)
	}
	mvn, _ := os.ReadFile(settings)
	mvnText := string(mvn)
	if !strings.Contains(mvnText, "<id>private</id>") || !strings.Contains(mvnText, "https://maven.aliyun.com/repository/public") {
		t.Fatalf("settings.xml = %s", mvnText)
	}
	if strings.Index(mvnText, "macpanel") > strings.Index(mvnText, "<id>private</id>") && strings.Contains(mvnText, "<mirrors>") {
		mirrorAt := strings.Index(mvnText, "<mirrors>")
		privateAt := strings.Index(mvnText, "<id>private</id>")
		if privateAt > mirrorAt && strings.Index(mvnText[mirrorAt:], "macpanel") > strings.Index(mvnText[mirrorAt:], "</mirrors>") {
			t.Fatal("macpanel mirror was not inserted inside mirrors")
		}
	}

	if err := Apply(opts, ApplyRequest{Ecosystem: "gradle", PresetID: "huawei"}); err != nil {
		t.Fatal(err)
	}
	gradlePath, _ := gradleInitPath(opts)
	gradleBody, err := os.ReadFile(gradlePath)
	if err != nil || !strings.Contains(string(gradleBody), "https://repo.huaweicloud.com/repository/maven/") {
		t.Fatalf("gradle script = %s %v", gradleBody, err)
	}

	goPath, _ := goEnvPath(opts)
	if err := os.MkdirAll(filepath.Dir(goPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goPath, []byte("GOPRIVATE=github.com/acme\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "go", PresetID: "goproxy-cn"}); err != nil {
		t.Fatal(err)
	}
	goBody, _ := os.ReadFile(goPath)
	if !strings.Contains(string(goBody), "GOPRIVATE=github.com/acme") || !strings.Contains(string(goBody), "GOPROXY=https://goproxy.cn,direct") || !strings.Contains(string(goBody), "GOSUMDB=sum.golang.google.cn") {
		t.Fatalf("go env = %s", goBody)
	}

	cargoPath, _ := cargoConfigPath(opts)
	if err := os.MkdirAll(filepath.Dir(cargoPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cargoPath, []byte("[net]\ngit-fetch-with-cli = true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "cargo", PresetID: "tsinghua"}); err != nil {
		t.Fatal(err)
	}
	cargoBody, _ := os.ReadFile(cargoPath)
	if !strings.Contains(string(cargoBody), "git-fetch-with-cli = true") || !strings.Contains(string(cargoBody), "sparse+https://mirrors.tuna.tsinghua.edu.cn/crates.io-index/") {
		t.Fatalf("cargo config = %s", cargoBody)
	}

	if err := os.WriteFile(opts.DockerDaemonPath, []byte("{\n  \"log-driver\": \"json-file\",\n  \"log-opts\": {\"max-size\": \"10m\"}\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "docker", PresetID: "daocloud"}); err != nil {
		t.Fatal(err)
	}
	daemon, _ := os.ReadFile(opts.DockerDaemonPath)
	daemonText := string(daemon)
	if !strings.Contains(daemonText, "\"log-driver\"") || !strings.Contains(daemonText, "https://docker.m.daocloud.io") {
		t.Fatalf("daemon.json = %s", daemonText)
	}
	if strings.Index(daemonText, "log-driver") > strings.Index(daemonText, "registry-mirrors") {
		t.Fatalf("daemon.json key order changed: %s", daemonText)
	}

	items, err := List(opts)
	if err != nil {
		t.Fatal(err)
	}
	active := map[string]string{}
	for _, item := range items {
		active[item.ID] = item.ActivePreset
		if item.ReadError != "" {
			t.Fatalf("%s read error: %s", item.ID, item.ReadError)
		}
		if item.Snippet == "" {
			t.Fatalf("%s missing snippet", item.ID)
		}
	}
	expects := map[string]string{
		"pip":    "aliyun",
		"npm":    "npmmirror",
		"maven":  "aliyun",
		"gradle": "huawei",
		"go":     "goproxy-cn",
		"cargo":  "tsinghua",
		"docker": "daocloud",
	}
	for id, preset := range expects {
		if active[id] != preset {
			t.Fatalf("%s active = %q, want %s", id, active[id], preset)
		}
	}
}

func TestOfficialClearsManagedValues(t *testing.T) {
	opts := testOpts(t)
	for _, eco := range []string{"pip", "npm", "maven", "gradle", "go", "cargo", "docker"} {
		preset := "tsinghua"
		switch eco {
		case "npm":
			preset = "npmmirror"
		case "maven", "gradle":
			preset = "aliyun"
		case "go":
			preset = "goproxy-cn"
		case "docker":
			preset = "daocloud"
		}
		if err := Apply(opts, ApplyRequest{Ecosystem: eco, PresetID: preset}); err != nil {
			t.Fatal(eco, err)
		}
		if err := Apply(opts, ApplyRequest{Ecosystem: eco, PresetID: "official"}); err != nil {
			t.Fatal(eco, err)
		}
	}
	items, err := List(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		if item.ActivePreset != "official" {
			t.Fatalf("%s active = %q after official", item.ID, item.ActivePreset)
		}
	}
	if _, err := os.Stat(filepath.Join(opts.Home, ".gradle", "init.d", "macpanel-mirror.init.gradle")); !os.IsNotExist(err) {
		t.Fatal("gradle init script still exists")
	}
}

func TestRejectsInvalidInputWithoutWriting(t *testing.T) {
	opts := testOpts(t)
	err := Apply(opts, ApplyRequest{
		Ecosystem: "npm",
		ValuesSet: true,
		Values:    map[string]string{"registry": "javascript:alert(1)"},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if _, statErr := os.Stat(filepath.Join(opts.Home, ".npmrc")); !os.IsNotExist(statErr) {
		t.Fatal("npmrc was written for an invalid registry")
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "pip", PresetID: "tsinghua", ValuesSet: true, Values: map[string]string{}}); err == nil {
		t.Fatal("expected mutual exclusion error")
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "missing", PresetID: "official"}); err == nil {
		t.Fatal("expected unknown ecosystem error")
	}
}

func TestDockerReadErrorDoesNotHideOtherEcosystems(t *testing.T) {
	opts := testOpts(t)
	if err := os.WriteFile(opts.DockerDaemonPath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := List(opts)
	if err != nil {
		t.Fatal(err)
	}
	var docker Ecosystem
	for _, item := range items {
		if item.ID == "docker" {
			docker = item
		}
		if item.ID == "pip" && item.ReadError != "" {
			t.Fatal(item.ReadError)
		}
	}
	if docker.ReadError == "" {
		t.Fatal("expected docker read error")
	}
	before, _ := os.ReadFile(opts.DockerDaemonPath)
	if err := Apply(opts, ApplyRequest{Ecosystem: "docker", PresetID: "daocloud"}); err == nil {
		t.Fatal("expected docker apply to fail on invalid json")
	}
	after, _ := os.ReadFile(opts.DockerDaemonPath)
	if string(before) != string(after) {
		t.Fatalf("invalid daemon.json changed: %s", after)
	}
}

func TestPathsFollowToolLocations(t *testing.T) {
	opts := testOpts(t)
	opts.GOOS = "darwin"
	pipPath, err := pipConfigPath(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(pipPath, filepath.Join("Library", "Application Support", "pip", "pip.conf")) {
		t.Fatal(pipPath)
	}
	goPath, err := goEnvPath(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(goPath, filepath.Join("Library", "Application Support", "go", "env")) {
		t.Fatal(goPath)
	}
	opts.GoEnvFile = filepath.Join(opts.Home, "custom-go-env")
	goPath, err = goEnvPath(opts)
	if err != nil || goPath != opts.GoEnvFile {
		t.Fatalf("go env path = %s err=%v", goPath, err)
	}
	opts.GOOS = "linux"
	opts.XDGConfigHome = filepath.Join(opts.Home, "xdg")
	opts.GoEnvFile = ""
	pipPath, err = pipConfigPath(opts)
	if err != nil || pipPath != filepath.Join(opts.XDGConfigHome, "pip", "pip.conf") {
		t.Fatalf("pip path = %s err=%v", pipPath, err)
	}
}

func TestCargoLeavesUnrelatedSource(t *testing.T) {
	opts := testOpts(t)
	path, _ := cargoConfigPath(opts)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	original := "[source.crates-io]\nreplace-with = \"ustc\"\n\n[source.ustc]\nregistry = \"sparse+https://mirrors.ustc.edu.cn/crates.io-index/\"\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := List(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		if item.ID == "cargo" && item.ActivePreset != "ustc" {
			t.Fatalf("cargo active = %s read=%s", item.ActivePreset, item.ReadError)
		}
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "cargo", PresetID: "official"}); err != nil {
		t.Fatal(err)
	}
	body, _ := os.ReadFile(path)
	if strings.Contains(string(body), "replace-with") {
		t.Fatalf("replace-with still present: %s", body)
	}
	if !strings.Contains(string(body), "[source.ustc]") {
		t.Fatalf("unrelated source removed: %s", body)
	}
}

func TestMavenOfficialKeepsServers(t *testing.T) {
	opts := testOpts(t)
	if err := Apply(opts, ApplyRequest{Ecosystem: "maven", PresetID: "tencent"}); err != nil {
		t.Fatal(err)
	}
	path, _ := mavenSettingsPath(opts)
	body, _ := os.ReadFile(path)
	updated := strings.Replace(string(body), "</settings>", "  <servers><server><id>keep</id></server></servers>\n</settings>", 1)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(opts, ApplyRequest{Ecosystem: "maven", PresetID: "official"}); err != nil {
		t.Fatal(err)
	}
	next, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(next), "macpanel") || !strings.Contains(string(next), "<id>keep</id>") {
		t.Fatalf("settings.xml = %s", next)
	}
}
