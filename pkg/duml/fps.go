package duml

type FPS int

const (
	UndefinedFPS = FPS(iota)
	FPS24
	FPS25
	FPS30
)

func (f FPS) BytesFixed() [1]byte {
	switch f {
	case FPS24:
		return [1]byte{0x01} // assumed, not confirmed
	case FPS25:
		return [1]byte{0x02}
	case FPS30:
		return [1]byte{0x03}
	default:
		return [1]byte{0}
	}
}

func FPSFromUint(v uint) FPS {
	switch v {
	case 24:
		return FPS24
	case 25:
		return FPS25
	case 30:
		return FPS30
	default:
		return UndefinedFPS
	}
}
