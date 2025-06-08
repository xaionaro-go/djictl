package djictl

type FPS int

const (
	FPSUndefined = FPS(iota)
	FPS25
	FPS30
)

func (f FPS) BytesFixed() [1]byte {
	switch f {
	case FPS25:
		return [1]byte{0x02}
	case FPS30:
		return [1]byte{0x03}
	default:
		return [1]byte{0}
	}
}
