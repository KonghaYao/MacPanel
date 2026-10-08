package terminal

import (
	"os"
	"strings"

	"github.com/1Panel-dev/1Panel/core/app/repo"
)

func InteractiveEnviron() []string {
	return applyInteractiveEnv(os.Environ(), utf8Locale())
}

func utf8Locale() string {
	lang, err := repo.NewISettingRepo().GetValueByKey("Language")
	if err != nil {
		return "en_US.UTF-8"
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(lang)), "zh") {
		return "zh_CN.UTF-8"
	}
	return "en_US.UTF-8"
}

func applyInteractiveEnv(base []string, locale string) []string {
	env := append([]string(nil), base...)
	env = setEnv(env, "TERM", "xterm-256color")
	env = setEnv(env, "COLORTERM", "truecolor")
	if envValue(env, "LANG") == "" {
		env = setEnv(env, "LANG", locale)
	}
	if envValue(env, "LC_ALL") == "" {
		env = setEnv(env, "LC_ALL", locale)
	}
	return env
}

func setEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			continue
		}
		out = append(out, item)
	}
	return append(out, prefix+value)
}

func envValue(env []string, key string) string {
	prefix := key + "="
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			return strings.TrimPrefix(item, prefix)
		}
	}
	return ""
}
