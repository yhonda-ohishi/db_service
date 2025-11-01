// +build windows

package main

import (
	"os"
	"syscall"
)

// getShutdownSignals returns the signals to listen for graceful shutdown on Windows
func getShutdownSignals() []os.Signal {
	// Windows only supports os.Interrupt (Ctrl+C, Ctrl+Break)
	// taskkill sends WM_CLOSE which may not be caught by signal.Notify
	return []os.Signal{os.Interrupt, syscall.SIGINT}
}
