package terminal

import "testing"

func TestApplyInteractiveEnvKeepsLocaleAndReplacesTerm(t *testing.T) {
	env := applyInteractiveEnv([]string{"TERM=xterm", "LANG=fr_FR.UTF-8", "LC_ALL=de_DE.UTF-8", "PATH=/bin"}, "zh_CN.UTF-8")
	if got := envValue(env, "TERM"); got != "xterm-256color" {
		t.Fatalf("TERM = %s", got)
	}
	if got := envValue(env, "COLORTERM"); got != "truecolor" {
		t.Fatalf("COLORTERM = %s", got)
	}
	if got := envValue(env, "LANG"); got != "fr_FR.UTF-8" {
		t.Fatalf("LANG = %s", got)
	}
	if got := envValue(env, "LC_ALL"); got != "de_DE.UTF-8" {
		t.Fatalf("LC_ALL = %s", got)
	}
	if countEnv(env, "TERM") != 1 {
		t.Fatalf("duplicate TERM: %v", env)
	}
}

func TestApplyInteractiveEnvFillsEmptyLocale(t *testing.T) {
	env := applyInteractiveEnv([]string{"PATH=/bin", "LANG=", "LC_ALL="}, "en_US.UTF-8")
	if got := envValue(env, "LANG"); got != "en_US.UTF-8" {
		t.Fatalf("LANG = %s", got)
	}
	if got := envValue(env, "LC_ALL"); got != "en_US.UTF-8" {
		t.Fatalf("LC_ALL = %s", got)
	}
}

func TestDockerExecArgs(t *testing.T) {
	args := dockerExecArgs([]string{"exec", "-e", "PGPASSWORD=secret", "-it", "db", "psql", "-U", "root"}, "zh_CN.UTF-8")
	joined := map[string]string{}
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-e" {
			key, value, ok := splitEnvArg(args[i+1])
			if ok {
				joined[key] = value
			}
		}
	}
	if joined["TERM"] != "xterm-256color" || joined["COLORTERM"] != "truecolor" {
		t.Fatalf("env flags = %#v in %v", joined, args)
	}
	if joined["PANEL_TERMINAL_LOCALE"] != "zh_CN.UTF-8" {
		t.Fatalf("locale flag = %#v", joined)
	}
	if _, overridden := joined["LANG"]; overridden {
		t.Fatalf("docker args override LANG: %v", args)
	}
	if args[len(args)-3] != "psql" || args[len(args)-1] != "root" {
		t.Fatalf("command was rewritten: %v", args)
	}
	foundPassword := false
	for _, arg := range args {
		if arg == "PGPASSWORD=secret" {
			foundPassword = true
		}
	}
	if !foundPassword {
		t.Fatalf("password arg missing: %v", args)
	}
}

func countEnv(env []string, key string) int {
	count := 0
	prefix := key + "="
	for _, item := range env {
		if len(item) >= len(prefix) && item[:len(prefix)] == prefix {
			count++
		}
	}
	return count
}

func splitEnvArg(arg string) (string, string, bool) {
	for i := 0; i < len(arg); i++ {
		if arg[i] == '=' {
			return arg[:i], arg[i+1:], true
		}
	}
	return "", "", false
}
