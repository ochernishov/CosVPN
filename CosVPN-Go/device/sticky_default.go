//go:build !linux

package device

import (
	"github.com/ochernishov/cosvpn/conn"
	"github.com/ochernishov/cosvpn/rwcancel"
)

func (device *Device) startRouteListener(_ conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
