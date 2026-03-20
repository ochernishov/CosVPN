/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package conn

import (
	"testing"
)

func TestPrettyName(t *testing.T) {
	var (
		recvFunc ReceiveFunc = func(bufs [][]byte, sizes []int, eps []Endpoint) (n int, err error) { return }
	)

	const want = "TestPrettyName"

	t.Run("ReceiveFunc.PrettyName", func(t *testing.T) {
		if got := recvFunc.PrettyName(); got != want {
			t.Errorf("PrettyName() = %v, want %v", got, want)
		}
	})
}
