package duml

import "encoding/binary"

const (
	RemoteControllerSimulatorStickMin    = uint16(364)
	RemoteControllerSimulatorStickCenter = uint16(1024)
	RemoteControllerSimulatorStickMax    = uint16(1684)

	RemoteControllerSimulatorDataSize = 38
)

// RemoteControllerSimulatorData represents the stick and button data sent in simulator mode.
// Reference: Remote Controller-N1 research.
type RemoteControllerSimulatorData struct {
	// Scaling: RemoteControllerSimulatorStickMin, RemoteControllerSimulatorStickCenter, RemoteControllerSimulatorStickMax
	RightStickHorizontal uint16
	RightStickVertical   uint16
	LeftStickVertical    uint16
	LeftStickHorizontal  uint16

	// Buttons (bitmask)
	Buttons uint32
}

func (d *RemoteControllerSimulatorData) Bytes() []byte {
	b := make([]byte, RemoteControllerSimulatorDataSize)
	binary.LittleEndian.PutUint16(b[0:2], d.RightStickHorizontal)
	binary.LittleEndian.PutUint16(b[2:4], d.RightStickVertical)
	binary.LittleEndian.PutUint16(b[4:6], d.LeftStickVertical)
	binary.LittleEndian.PutUint16(b[6:8], d.LeftStickHorizontal)
	binary.LittleEndian.PutUint32(b[8:12], d.Buttons)
	// Fill rest with 0 or as per protocol
	return b
}

func NewRemoteControllerSimulatorMessage(data RemoteControllerSimulatorData) *Message {
	return &Message{
		Interface: InterfaceIDAppToRemoteController,
		ID:        MessageID(0), // Initial sequence number
		Type:      MessageTypeRemoteControllerSimulatorData,
		Payload:   data.Bytes(),
	}
}
