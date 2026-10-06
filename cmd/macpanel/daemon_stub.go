//go:build !darwin

package main

import "fmt"

func isDaemonChild() bool {
	return false
}

func startDaemon() error {
	return fmt.Errorf("daemon mode is only supported on macOS")
}

func writeDaemonPID() {}

func removeDaemonPID() {}
