package duml

func NewFCCEnableMessage(value bool) *Message {
	payload := uint8(0x00)
	if value {
		payload = uint8(0x01)
	}
	return &Message{
		Type:    MessageTypeFCCSupport,
		Payload: []byte{payload}, // Simple enable flag
	}
}
