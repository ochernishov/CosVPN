/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
 */

package device

// DisableSomeRoamingForBrokenMobileSemantics should ideally be called before peers are created,
// though it will try to deal with it, and race maybe, if called after.
func (device *Device) DisableSomeRoamingForBrokenMobileSemantics() {
	device.net.brokenRoaming = true
	device.peers.RLock()
	for _, peer := range device.peers.keyMap {
		peer.endpoint.Lock()
		peer.endpoint.disableRoaming = peer.endpoint.val != nil
		peer.endpoint.Unlock()
	}
	device.peers.RUnlock()
}
