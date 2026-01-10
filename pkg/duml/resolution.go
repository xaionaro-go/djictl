package duml

import "strings"

type Resolution int

const (
	UndefinedResolution = Resolution(iota)
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

func ResolutionFromString(s string) Resolution {
	switch strings.ToLower(s) {
	case "480p", "640x480":
		return Resolution480p
	case "720p", "1280x720":
		return Resolution720p
	case "1080p", "1920x1080":
		return Resolution1080p
	default:
		return UndefinedResolution
	}
}
