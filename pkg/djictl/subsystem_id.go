package djictl

import "fmt"

type SubsystemID uint16

const (
	SubsystemIDConfigurer = SubsystemID(0x0201)
	SubsystemIDPairer     = SubsystemID(0x0207)
	SubsystemIDStreamer   = SubsystemID(0x0208)

	SubsystemIDPrePairer = SubsystemID(0x0402)

	SubsystemIDOneMorePairer = SubsystemID(0x0288)
)

func (id SubsystemID) String() string {
	switch id {
	case SubsystemIDConfigurer:
		return "configurer"
	case SubsystemIDPrePairer:
		return "pre-pairer"
	case SubsystemIDPairer:
		return "pairer"
	case SubsystemIDStreamer:
		return "streamer"
	default:
		return fmt.Sprintf("%04X", uint16(id))
	}
}
