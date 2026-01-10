package duml

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseMessage_Padding(t *testing.T) {
	// Sample from wlan0.pcap (extracted via tshark)
	// 55 11 04 92 02 1b fd 94 40 07 45 00 00 00 00 31 d7
	rawHex := "55110492021bfd944007450000000031d7"
	raw, _ := hex.DecodeString(rawHex)

	t.Run("no padding", func(t *testing.T) {
		msg, err := ParseMessage(raw)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
	})

	t.Run("with padding", func(t *testing.T) {
		padded := append(raw, 1, 2, 3, 4) // non-zero padding
		msg, err := ParseMessage(padded)
		assert.NoError(t, err, "ParseMessage should succeed even with padding")
		assert.NotNil(t, msg)
	})
}

func TestParseMessage_10BitLength(t *testing.T) {
	msg := &Message{
		Interface: InterfaceID{Sender: 1, Receiver: 2},
		ID:        3,
		Type:      MessageTypeGetVersion,
		Payload:   make([]byte, 300), // Resulting total length will be 313 (300+13)
	}
	for i := range msg.Payload {
		msg.Payload[i] = byte(i)
	}

	raw := msg.Bytes()
	assert.Equal(t, 313, len(raw))

	parsed, err := ParseMessage(raw)
	require.NoError(t, err)
	assert.Equal(t, msg.Payload, parsed.Payload)
}

func TestPcapSamples(t *testing.T) {
	samples := []string{
		"551204c70402f6010004270000080000299d",                           // frame 645
		"550d04330207ea94400707242b",                                     // frame 747
		"551f044e0702ea94c0070700104f736d6f506f636b6574332d36303934ccc8", // frame 753
		"553e044b0402f001000405010700000b008000e0ff00014987fc276afd00001bfd010013aed43b68ce6c3d00917f3ff91bc5b900000000010000000099b5", // frame 643
	}

	for _, s := range samples {
		t.Run(s, func(t *testing.T) {
			b, err := hex.DecodeString(s)
			require.NoError(t, err)
			msg, err := ParseMessage(b)
			if err != nil {
				t.Fatalf("failed to parse message %s: %v", s, err)
			}
			require.NotNil(t, msg)
			t.Logf("CRC16 for %s: %04X", s, crc16(b[:len(b)-2]))
		})
	}
}
