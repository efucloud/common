// Copyright (c) citicbank.com All rights reserved
//go:build windows
// +build windows

package signals

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
