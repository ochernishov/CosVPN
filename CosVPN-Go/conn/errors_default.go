//go:build !linux

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package conn

func errShouldDisableUDPGSO(_ error) bool {
	return false
}
