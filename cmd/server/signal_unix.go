// +build !windows

package main

import (
	"os"
	"syscall"
)

// getShutdownSignals returns the signals to listen for graceful shutdown on Unix systems
func getShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
}
