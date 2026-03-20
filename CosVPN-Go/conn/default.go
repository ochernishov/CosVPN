//go:build !windows

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package conn

func NewDefaultBind() Bind { return NewStdNetBind() }
