package djictl

type Resolution int

const (
	ResolutionUndefined = Resolution(iota)
	Resolution480p
	Resolution720p
	Resolution1080p
)

func (r Resolution) BytesFixed() [1]byte {
	switch r {
	case Resolution480p:
		return [1]byte{0x47}
	case Resolution720p:
		return [1]byte{0x04}
	case Resolution1080p:
		return [1]byte{0x0A}
	default:
		return [1]byte{0}
	}
}
