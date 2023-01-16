// Copyright (c) citicbank.com All rights reserved
//go:build linux || darwin
// +build linux darwin

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
