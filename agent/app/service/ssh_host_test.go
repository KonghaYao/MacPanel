package service

import "testing"

func TestValidateSSHHostInput(t *testing.T) {
	if err := validateSSHHostInput("my-server", "192.168.1.10", "ubuntu"); err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	if err := validateSSHHostInput("bad alias", "192.168.1.10", "ubuntu"); err == nil {
		t.Fatal("expected alias validation error")
	}
	if err := validateSSHHostInput("server", "not a host!", "ubuntu"); err == nil {
		t.Fatal("expected hostname validation error")
	}
	if err := validateSSHHostInput("server", "example.com", "bad user!"); err == nil {
		t.Fatal("expected user validation error")
	}
}

func TestIsValidSSHHostName(t *testing.T) {
	cases := map[string]bool{
		"192.168.0.1":   true,
		"example.com":   true,
		"[::1]":         true,
		"bad host name": false,
		"":              false,
	}
	for host, expected := range cases {
		if got := isValidSSHHostName(host); got != expected {
			t.Fatalf("host %q expected %v got %v", host, expected, got)
		}
	}
}
