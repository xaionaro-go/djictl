package duml

type ImageStabilization int

const (
	ImageStabilizationUndefined = ImageStabilization(iota)
	ImageStabilizationOff
	ImageStabilizationRockSteady
	ImageStabilizationRockSteadyPlus
	ImageStabilizationHorizonBalancing
	ImageStabilizationHorizonSteady
)

func (v ImageStabilization) BytesFixed() [1]byte {
	switch v {
	case ImageStabilizationOff:
		return [1]byte{0x00}
	case ImageStabilizationRockSteady:
		return [1]byte{0x01}
	case ImageStabilizationHorizonSteady:
		return [1]byte{0x02}
	case ImageStabilizationRockSteadyPlus:
		return [1]byte{0x03}
	case ImageStabilizationHorizonBalancing:
		return [1]byte{0x04}
	}
	return [1]byte{0x00}
}
