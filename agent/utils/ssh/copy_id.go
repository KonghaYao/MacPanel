package ssh

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

func CopyPublicKeyWithPassword(host, user, password string, port int, publicKey []byte) error {
	if port == 0 {
		port = 22
	}
	pubKeyLine := strings.TrimSpace(string(publicKey))
	if pubKeyLine == "" {
		return fmt.Errorf("public key is empty")
	}

	config := &gossh.ClientConfig{
		User:            user,
		Auth:            []gossh.AuthMethod{gossh.Password(password)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	client, err := gossh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}
	defer client.Close()

	keyB64 := base64.StdEncoding.EncodeToString([]byte(pubKeyLine))
	script := fmt.Sprintf(
		"mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && KEY=$(echo %s | base64 -d) && grep -qxF \"$KEY\" ~/.ssh/authorized_keys 2>/dev/null || printf '%%s\\n' \"$KEY\" >> ~/.ssh/authorized_keys",
		shellQuote(keyB64),
	)

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("create ssh session failed: %w", err)
	}
	defer session.Close()

	if err := session.Run(script); err != nil {
		return fmt.Errorf("install public key failed: %w", err)
	}
	return nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
