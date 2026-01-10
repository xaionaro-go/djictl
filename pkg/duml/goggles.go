package duml

// Goggles mode switching
const ()

type GogglesMode uint8

const (
	GogglesModeNormal GogglesMode = 0
	GogglesModeUSB    GogglesMode = 1
)

func NewGogglesModeMessage(mode GogglesMode) *Message {
	return &Message{
		Type:    MessageTypeGogglesMode,
		Payload: []byte{uint8(mode)},
	}
}
