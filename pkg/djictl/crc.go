package djictl

import (
	"fmt"

	crc16upstream "github.com/sigurn/crc16"
	crc8upstream "github.com/sigurn/crc8"
)

var crc8Table = crc8upstream.MakeTable(crc8upstream.Params{
	Poly:   0x31,
	Init:   0xEE,
	RefIn:  true,
	RefOut: true,
	XorOut: 0x00,
	Name:   "DJI/CRC8",
})

var crc16Table = crc16upstream.MakeTable(crc16upstream.Params{
	Poly:   0x1021,
	Init:   0x496C,
	RefIn:  true,
	RefOut: true,
	XorOut: 0x00,
	Name:   "DJI/CRC16",
})

func init() {
	if v := crc8upstream.Init(crc8Table); v != 0xEE {
		panic(fmt.Errorf("internal error, the initial value is supposed to be 0xEE, but received 0x%02X instead", v))
	}

	if v := crc16upstream.Init(crc16Table); v != 0x496C {
		panic(fmt.Errorf("internal error, the initial value is supposed to be 0x496C, but received 0x%04X instead", v))
	}
}

func crc8(in []byte) uint8 {
	return crc8upstream.Checksum(in, crc8Table)
}

func crc16(in []byte) uint16 {
	return crc16upstream.Checksum(in, crc16Table)
}
