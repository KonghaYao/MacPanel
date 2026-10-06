//go:build darwin

package websocket

import "testing"

func TestParseDarwinWhoLineRemote(t *testing.T) {
	session, ok := parseDarwinWhoLine("john    ttys001      Oct  6 10:00 00:01    12345 (192.168.1.5)")
	if !ok {
		t.Fatal("expected remote who line to parse")
	}
	if session.Username != "john" || session.Host != "192.168.1.5" || session.PID != 12345 {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestParseDarwinWhoLineLocal(t *testing.T) {
	_, ok := parseDarwinWhoLine("mino             ttys006      Oct  6 09:58 00:23    84669")
	if ok {
		t.Fatal("expected local who line without host to be skipped")
	}
}
