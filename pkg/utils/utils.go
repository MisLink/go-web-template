package utils

import (
	"os"
	"syscall"
)

var GracefulShutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT}
