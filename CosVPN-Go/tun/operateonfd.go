//go:build darwin || freebsd

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package tun

import (
	"fmt"
)

func (tun *NativeTun) operateOnFd(fn func(fd uintptr)) {
	sysconn, err := tun.tunFile.SyscallConn()
	if err != nil {
		tun.errors <- fmt.Errorf("unable to find sysconn for tunfile: %s", err.Error())
		return
	}
	err = sysconn.Control(fn)
	if err != nil {
		tun.errors <- fmt.Errorf("unable to control sysconn for tunfile: %s", err.Error())
	}
}
