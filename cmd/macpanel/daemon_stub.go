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

func stopDaemon() error {
	return fmt.Errorf("daemon control is only supported on macOS")
}

func restartDaemon() error {
	return fmt.Errorf("daemon control is only supported on macOS")
}

func printDaemonStatus() error {
	return fmt.Errorf("daemon control is only supported on macOS")
}
