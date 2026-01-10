package duml

// FCC Mode Trigger (Product Shielded Config)
const (
	FCCEnablePayload = uint8(0x01)
)

func NewFCCEnableMessage() *Message {
	return &Message{
		Type:    MessageTypeFCCSupport,
		Payload: []byte{FCCEnablePayload}, // Simple enable flag
	}
}
