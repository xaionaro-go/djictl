package duml

import (
	"testing"
)

func TestGogglesModeMessage(t *testing.T) {
	msg := NewGogglesModeMessage(GogglesModeUSB)
	if msg.Type != MessageTypeGogglesMode {
		t.Errorf("Expected type %v, got %v", MessageTypeGogglesMode, msg.Type)
	}
	if msg.Payload[0] != 1 {
		t.Errorf("Expected payload 1, got %v", msg.Payload[0])
	}
	// Verify serialization doesn't panic
	b := msg.Bytes()
	if len(b) == 0 {
		t.Errorf("Serialized message is empty")
	}
}

func TestFCCEnableMessage(t *testing.T) {
	msg := NewFCCEnableMessage(true)
	if msg.Type != MessageTypeFCCSupport {
		t.Errorf("Expected type %v, got %v", MessageTypeFCCSupport, msg.Type)
	}
	b := msg.Bytes()
	if len(b) == 0 {
		t.Errorf("Serialized message is empty")
	}
}

func TestRemoteControllerSimulatorMessage(t *testing.T) {
	data := RemoteControllerSimulatorData{
		RightStickHorizontal: 1024,
		RightStickVertical:   1024,
		LeftStickVertical:    1024,
		LeftStickHorizontal:  1024,
	}
	msg := NewRemoteControllerSimulatorMessage(data)
	if msg.Type != MessageTypeRemoteControllerSimulatorData {
		t.Errorf("Expected type %v, got %v", MessageTypeRemoteControllerSimulatorData, msg.Type)
	}
	b := msg.Bytes()
	if len(b) != 13+38 { // 13 bytes header/overhead + 38 bytes payload
		t.Errorf("Expected length 51, got %d", len(b))
	}
}
