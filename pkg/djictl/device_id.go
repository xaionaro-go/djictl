package djictl

import (
	"net"
)

type DeviceID = net.HardwareAddr

func ParseDeviceID(s string) (DeviceID, error) {
	return net.ParseMAC(s)
}
