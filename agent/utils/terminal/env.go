package terminal

import (
	"os"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/repo"
)

func InteractiveEnviron() []string {
	return applyInteractiveEnv(os.Environ(), utf8Locale())
}

func DockerExecArgs(args []string) []string {
	return dockerExecArgs(args, utf8Locale())
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

func dockerExecArgs(args []string, locale string) []string {
	if len(args) == 0 || args[0] != "exec" {
		return args
	}
	index := 1
	var opts []string
	for index < len(args) {
		arg := args[index]
		if arg == "--" {
			index++
			break
		}
		if !strings.HasPrefix(arg, "-") {
			break
		}
		opts = append(opts, arg)
		if dockerOptionTakesValue(arg) && index+1 < len(args) {
			index++
			opts = append(opts, args[index])
		}
		index++
	}
	if index >= len(args) {
		return args
	}
	container := args[index]
	command := args[index+1:]
	if len(command) == 0 {
		command = []string{"sh"}
	}
	script := `if [ -z "$LANG" ]; then export LANG="$PANEL_TERMINAL_LOCALE"; fi; if [ -z "$LC_ALL" ]; then export LC_ALL="$PANEL_TERMINAL_LOCALE"; fi; exec "$@"`
	wrapped := append([]string{"sh", "-c", script, "sh"}, command...)
	out := append([]string{"exec"}, opts...)
	out = append(out,
		"-e", "TERM=xterm-256color",
		"-e", "COLORTERM=truecolor",
		"-e", "PANEL_TERMINAL_LOCALE="+locale,
		container,
	)
	return append(out, wrapped...)
}

func dockerOptionTakesValue(arg string) bool {
	switch arg {
	case "-u", "--user", "-e", "--env", "-w", "--workdir":
		return true
	default:
		return false
	}
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
