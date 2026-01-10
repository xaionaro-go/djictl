package duml

import (
	"encoding/hex"
	"testing"
)

func TestPacketPcap_CRC8(t *testing.T) {
	// Sample from app_to_drone.txt: 55 1b 04 75
	data, _ := hex.DecodeString("551b04")
	expected := uint8(0x75)
	actual := crc8(data)
	if actual != expected {
		t.Errorf("CRC8 mismatch: expected %02X, got %02X", expected, actual)
	}
}

func TestPacketPcap_CRC8_V2(t *testing.T) {
	// Sample from app_to_drone.txt: 55 12 04 c7
	data, _ := hex.DecodeString("551204")
	expected := uint8(0xc7)
	actual := crc8(data)
	if actual != expected {
		t.Errorf("CRC8 mismatch for v2: expected %02X, got %02X", expected, actual)
	}
}
