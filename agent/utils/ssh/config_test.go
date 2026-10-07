package ssh

import (
	"strings"
	"testing"
)

func TestParseSSHConfigPreservesForeignBlocks(t *testing.T) {
	content := `# user config
Host github.com
  HostName github.com
  User git

# Managed by MacPanel
Host myserver
  HostName 192.168.1.10
  User admin
`
	blocks := ParseSSHConfig(content)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	if blocks[0].Alias != "github.com" || blocks[0].Managed {
		t.Fatalf("unexpected first block: %+v", blocks[0])
	}
	if blocks[1].Alias != "myserver" || !blocks[1].Managed {
		t.Fatalf("unexpected second block: %+v", blocks[1])
	}
	if blocks[1].HostName != "192.168.1.10" || blocks[1].User != "admin" {
		t.Fatalf("unexpected host fields: %+v", blocks[1])
	}
}

func TestAddAndRemoveSSHConfigHost(t *testing.T) {
	base := `Host github.com
  HostName github.com
  User git
`
	updated, err := AddSSHConfigHost(base, "lab", "10.0.0.5", "ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	if !HasSSHConfigHost(updated, "lab") {
		t.Fatal("expected lab host to exist")
	}
	if !HasSSHConfigHost(updated, "github.com") {
		t.Fatal("expected github host to remain")
	}

	removed, err := RemoveSSHConfigHost(updated, "lab")
	if err != nil {
		t.Fatal(err)
	}
	if HasSSHConfigHost(removed, "lab") {
		t.Fatal("expected lab host to be removed")
	}
	if !HasSSHConfigHost(removed, "github.com") {
		t.Fatal("expected github host to remain after removal")
	}
	if strings.Contains(removed, managedMarker+"\nHost lab") {
		t.Fatal("managed block should be fully removed")
	}
}

func TestAddSSHConfigHostRejectsDuplicate(t *testing.T) {
	content := FormatSSHConfigBlock("lab", "10.0.0.5", "ubuntu")
	if _, err := AddSSHConfigHost(content, "lab", "10.0.0.6", "root"); err == nil {
		t.Fatal("expected duplicate alias error")
	}
}
