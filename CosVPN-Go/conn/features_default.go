//go:build !linux
// +build !linux

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package conn

import "net"

func supportsUDPOffload(_ *net.UDPConn) (txOffload, rxOffload bool) {
	return
}
