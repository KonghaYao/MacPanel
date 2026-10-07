package mirrors

import "strings"

func definitions() []definition {
	return []definition{
		pipDefinition(),
		npmDefinition(),
		mavenDefinition(),
		gradleDefinition(),
		goDefinition(),
		cargoDefinition(),
		dockerDefinition(),
	}
}

func pipDefinition() definition {
	return definition{
		id: "pip",
		fields: []Field{
			{Key: "indexUrl", Kind: KindURL},
			{Key: "trustedHost", Kind: KindHost},
		},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "tsinghua", name: "Tsinghua", values: map[string]string{
				"indexUrl":    "https://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple",
				"trustedHost": "mirrors.tuna.tsinghua.edu.cn",
			}},
			{id: "aliyun", name: "Aliyun", values: map[string]string{
				"indexUrl":    "https://mirrors.aliyun.com/pypi/simple/",
				"trustedHost": "mirrors.aliyun.com",
			}},
			{id: "ustc", name: "USTC", values: map[string]string{
				"indexUrl":    "https://pypi.mirrors.ustc.edu.cn/simple",
				"trustedHost": "pypi.mirrors.ustc.edu.cn",
			}},
			{id: "tencent", name: "Tencent Cloud", values: map[string]string{
				"indexUrl":    "https://mirrors.cloud.tencent.com/pypi/simple",
				"trustedHost": "mirrors.cloud.tencent.com",
			}},
		},
		configPath: pipConfigPath,
		read:       readPip,
		write:      writePip,
		snippet:    pipSnippet,
	}
}

func npmDefinition() definition {
	return definition{
		id:     "npm",
		fields: []Field{{Key: "registry", Kind: KindURL}},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "npmmirror", name: "npmmirror", values: map[string]string{"registry": "https://registry.npmmirror.com"}},
			{id: "huawei", name: "Huawei Cloud", values: map[string]string{"registry": "https://repo.huaweicloud.com/repository/npm/"}},
			{id: "tencent", name: "Tencent Cloud", values: map[string]string{"registry": "https://mirrors.cloud.tencent.com/npm/"}},
		},
		configPath: npmrcPath,
		read:       readNpm,
		write:      writeNpm,
		snippet:    npmSnippet,
	}
}

func mavenDefinition() definition {
	return definition{
		id:         "maven",
		fields:     []Field{{Key: "repository", Kind: KindURL}},
		presets:    mavenPresets(),
		configPath: mavenSettingsPath,
		read:       readMaven,
		write:      writeMaven,
		snippet:    mavenSnippet,
	}
}

func gradleDefinition() definition {
	return definition{
		id:         "gradle",
		fields:     []Field{{Key: "repository", Kind: KindURL}},
		presets:    mavenPresets(),
		configPath: gradleInitPath,
		read:       readGradle,
		write:      writeGradle,
		snippet:    gradleSnippet,
	}
}

func mavenPresets() []presetDef {
	return []presetDef{
		{id: "official", name: "Official", values: map[string]string{}},
		{id: "aliyun", name: "Aliyun", values: map[string]string{"repository": "https://maven.aliyun.com/repository/public"}},
		{id: "huawei", name: "Huawei Cloud", values: map[string]string{"repository": "https://repo.huaweicloud.com/repository/maven/"}},
		{id: "tencent", name: "Tencent Cloud", values: map[string]string{"repository": "https://mirrors.cloud.tencent.com/nexus/repository/maven-public/"}},
	}
}

func goDefinition() definition {
	return definition{
		id: "go",
		fields: []Field{
			{Key: "goproxy", Kind: KindGoproxy},
			{Key: "gosumdb", Kind: KindGosumdb},
		},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "goproxy-cn", name: "goproxy.cn", values: map[string]string{
				"goproxy": "https://goproxy.cn,direct",
				"gosumdb": "sum.golang.google.cn",
			}},
			{id: "aliyun", name: "Aliyun", values: map[string]string{
				"goproxy": "https://mirrors.aliyun.com/goproxy/,direct",
				"gosumdb": "sum.golang.google.cn",
			}},
			{id: "goproxy-io", name: "goproxy.io", values: map[string]string{
				"goproxy": "https://goproxy.io,direct",
				"gosumdb": "sum.golang.google.cn",
			}},
		},
		configPath: goEnvPath,
		read:       readGo,
		write:      writeGo,
		snippet:    goSnippet,
	}
}

func cargoDefinition() definition {
	return definition{
		id:     "cargo",
		fields: []Field{{Key: "registry", Kind: KindCargo}},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "tsinghua", name: "Tsinghua", values: map[string]string{"registry": "sparse+https://mirrors.tuna.tsinghua.edu.cn/crates.io-index/"}},
			{id: "ustc", name: "USTC", values: map[string]string{"registry": "sparse+https://mirrors.ustc.edu.cn/crates.io-index/"}},
			{id: "rsproxy", name: "rsproxy", values: map[string]string{"registry": "sparse+https://rsproxy.cn/index/"}},
		},
		configPath: cargoConfigPath,
		read:       readCargo,
		write:      writeCargo,
		snippet:    cargoSnippet,
	}
}

func dockerDefinition() definition {
	return definition{
		id:     "docker",
		fields: []Field{{Key: "mirrors", Kind: KindURLs}},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "daocloud", name: "DaoCloud", values: map[string]string{"mirrors": "https://docker.m.daocloud.io"}},
			{id: "panel", name: "1Panel", values: map[string]string{"mirrors": "https://docker.1panel.live"}},
		},
		configPath: dockerDaemonPath,
		read:       readDocker,
		write:      writeDocker,
		snippet:    dockerSnippet,
	}
}

func homebrewDefinition() definition {
	return definition{
		id: EcosystemHomebrew,
		fields: []Field{
			{Key: "bottleDomain", Kind: KindURL},
			{Key: "apiDomain", Kind: KindURL},
			{Key: "brewGitRemote", Kind: KindURL},
			{Key: "coreGitRemote", Kind: KindURL},
			{Key: "caskGitRemote", Kind: KindURL},
		},
		presets: []presetDef{
			{id: "official", name: "Official", values: map[string]string{}},
			{id: "tsinghua", name: "Tsinghua", values: map[string]string{
				"bottleDomain":  "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles",
				"apiDomain":     "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api",
				"brewGitRemote": "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git",
				"coreGitRemote": "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git",
				"caskGitRemote": "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-cask.git",
			}},
			{id: "ustc", name: "USTC", values: map[string]string{
				"bottleDomain":  "https://mirrors.ustc.edu.cn/homebrew-bottles",
				"apiDomain":     "https://mirrors.ustc.edu.cn/homebrew-bottles/api",
				"brewGitRemote": "https://mirrors.ustc.edu.cn/brew.git",
				"coreGitRemote": "https://mirrors.ustc.edu.cn/homebrew-core.git",
				"caskGitRemote": "https://mirrors.ustc.edu.cn/homebrew-cask.git",
			}},
			{id: "aliyun", name: "Aliyun", values: map[string]string{
				"bottleDomain":  "https://mirrors.aliyun.com/homebrew/homebrew-bottles",
				"apiDomain":     "https://mirrors.aliyun.com/homebrew/homebrew-bottles/api",
				"brewGitRemote": "https://mirrors.aliyun.com/homebrew/brew.git",
				"coreGitRemote": "https://mirrors.aliyun.com/homebrew/homebrew-core.git",
				"caskGitRemote": "https://mirrors.aliyun.com/homebrew/homebrew-cask.git",
			}},
		},
		snippet: homebrewSnippet,
	}
}

func HomebrewState(current map[string]string, configPath string) Ecosystem {
	return BuildState(homebrewDefinition(), current, configPath, nil)
}

func ResolveHomebrew(req ApplyRequest) (map[string]string, error) {
	values, err := ResolveValues(homebrewDefinition(), req)
	if err != nil {
		return nil, err
	}
	if err := validateValues(homebrewDefinition(), values); err != nil {
		return nil, err
	}
	return values, nil
}

func pipSnippet(values map[string]string) string {
	if values["indexUrl"] == "" && values["trustedHost"] == "" {
		return "pip config unset global.index-url\npip config unset global.trusted-host"
	}
	var lines []string
	if values["indexUrl"] != "" {
		lines = append(lines, "pip config set global.index-url "+values["indexUrl"])
	}
	if values["trustedHost"] != "" {
		lines = append(lines, "pip config set global.trusted-host "+values["trustedHost"])
	}
	return strings.Join(lines, "\n")
}

func npmSnippet(values map[string]string) string {
	if values["registry"] == "" {
		return "npm config delete registry"
	}
	return "npm config set registry " + values["registry"]
}

func mavenSnippet(values map[string]string) string {
	if values["repository"] == "" {
		return "<!-- remove the mirror whose id is macpanel from ~/.m2/settings.xml -->"
	}
	return "<mirror>\n  <id>macpanel</id>\n  <url>" + values["repository"] + "</url>\n  <mirrorOf>*</mirrorOf>\n</mirror>"
}

func gradleSnippet(values map[string]string) string {
	if values["repository"] == "" {
		return "# delete ~/.gradle/init.d/macpanel-mirror.init.gradle"
	}
	return strings.Replace(gradleScriptTemplate, "__MIRROR__", values["repository"], 1)
}

func goSnippet(values map[string]string) string {
	if values["goproxy"] == "" && values["gosumdb"] == "" {
		return "go env -u GOPROXY\ngo env -u GOSUMDB"
	}
	var lines []string
	if values["goproxy"] != "" {
		lines = append(lines, "go env -w GOPROXY="+values["goproxy"])
	}
	if values["gosumdb"] != "" {
		lines = append(lines, "go env -w GOSUMDB="+values["gosumdb"])
	}
	return strings.Join(lines, "\n")
}

func cargoSnippet(values map[string]string) string {
	if values["registry"] == "" {
		return "# remove [source.macpanel], [registries.macpanel], and source.crates-io replace-with from $CARGO_HOME/config.toml"
	}
	return "[source.crates-io]\nreplace-with = \"macpanel\"\n\n[source.macpanel]\nregistry = \"" + values["registry"] + "\"\n\n[registries.macpanel]\nindex = \"" + values["registry"] + "\""
}

func dockerSnippet(values map[string]string) string {
	mirrors := splitURLList(values["mirrors"])
	if len(mirrors) == 0 {
		return "# delete registry-mirrors from daemon.json"
	}
	var b strings.Builder
	b.WriteString("{\n  \"registry-mirrors\": [\n")
	for i, mirror := range mirrors {
		b.WriteString("    \"" + mirror + "\"")
		if i != len(mirrors)-1 {
			b.WriteString(",")
		}
		b.WriteByte('\n')
	}
	b.WriteString("  ]\n}")
	return b.String()
}

func homebrewSnippet(values map[string]string) string {
	pairs := []struct {
		key string
		env string
	}{
		{"bottleDomain", "HOMEBREW_BOTTLE_DOMAIN"},
		{"apiDomain", "HOMEBREW_API_DOMAIN"},
		{"brewGitRemote", "HOMEBREW_BREW_GIT_REMOTE"},
		{"coreGitRemote", "HOMEBREW_CORE_GIT_REMOTE"},
		{"caskGitRemote", "HOMEBREW_CASK_GIT_REMOTE"},
	}
	var lines []string
	for _, pair := range pairs {
		if values[pair.key] != "" {
			lines = append(lines, pair.env+"="+values[pair.key])
		}
	}
	if len(lines) == 0 {
		return "# MacPanel does not inject HOMEBREW_* mirror variables"
	}
	return strings.Join(lines, "\n")
}
