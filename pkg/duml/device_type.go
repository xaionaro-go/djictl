package duml

import "bytes"

type DeviceType int

const (
	DeviceTypeUndefined = DeviceType(iota)
	DeviceTypeUnknown
	DeviceTypeOsmoAction3
	DeviceTypeOsmoAction4
	DeviceTypeOsmoAction5Pro
	DeviceTypeOsmoPocket3
	DeviceTypeMiniSE
	DeviceTypeAir2S
	DeviceTypeMavic3
	EndOfDeviceType
)

func (t DeviceType) Magic() [2]byte {
	switch t {
	case DeviceTypeOsmoAction3:
		return [2]byte{0x12, 0x00}
	case DeviceTypeOsmoAction4:
		return [2]byte{0x14, 0x00}
	case DeviceTypeOsmoAction5Pro:
		return [2]byte{0x15, 0x00}
	case DeviceTypeOsmoPocket3:
		return [2]byte{0x20, 0x00}
	case DeviceTypeAir2S:
		return [2]byte{0x17, 0x00}
	case DeviceTypeMiniSE:
		return [2]byte{0x19, 0x00}
	case DeviceTypeMavic3:
		return [2]byte{0x1C, 0x00}
	}
	return [2]byte{0, 0}
}

func (t DeviceType) BytesFixedStartStreaming() [1]byte {
	switch t {
	case DeviceTypeOsmoAction5Pro:
		return [1]byte{0x2E}
	}
	return [1]byte{0x2A}
}

func (t DeviceType) BytesFixedSetImageStabilization() [1]byte {
	switch t {
	case DeviceTypeOsmoAction5Pro:
		return [1]byte{0x1A}
	}
	return [1]byte{0x08}
}

var djiMagic = []byte{0xAA, 0x08}

func IdentifyDeviceType(manufacturerData []byte) DeviceType {
	if !bytes.HasPrefix(manufacturerData, djiMagic) {
		return DeviceTypeUndefined
	}

	for t := DeviceTypeUndefined + 2; t < EndOfDeviceType; t++ {
		magic := t.Magic()
		if bytes.Equal(manufacturerData[2:4], magic[:]) {
			return t
		}
	}
	return DeviceTypeUnknown
}
