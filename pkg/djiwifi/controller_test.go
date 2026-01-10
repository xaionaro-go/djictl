package djiwifi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func TestPacket_Serialization(t *testing.T) {
	// Frame 636 from wlan0.pcap
	// Length: 48 (0x30)
	// Type: 0x80
	// Signature: 0x47 0xa8
	// Metadata: 00 00 00 5f 38 42 64 00 64 00 c0 05 14 00 00 64
	// Payload: 00 00 01 90 01 c0 05 14 00 00 64 00 14 00 64 00 c0 05 14 00 00 64 00 01 01 04 01 02
	rawHex := "308047a80000005f384264006400c005140000640000019001c005140000640014006400c00514000064000101040102"
	raw, _ := hex.DecodeString(rawHex)

	p, err := ParsePacket(raw)
	assert.NoError(t, err)
	if p == nil {
		t.Fatal("packet is nil")
	}
	assert.Equal(t, MessageTypeStandard, p.Type)
	assert.Equal(t, raw, p.Bytes())
}

func TestPacket_DUML(t *testing.T) {
	// Frame 702 from wlan0.pcap
	// Wrapper: 4b 80 47 a8 e8 42 05 8b 68 42 e8 42 00 00 00 00 16 01 00 00
	// DUML: 55 37 04 f9 02 28 de 94 40 00 99 02 02 00 00 4b 34 00 00 00 00 00 1d 00 17 00 70 72 6f 64 75 63 74 5f 73 68 69 65 6c 64 65 64 5f 63 6f 6e 66 69 67 00 00 00 00 00 66
	rawHex := "4b8047a8e842058b6842e8420000000016010000553704f90228de94400099020200004b3400000000001d00170070726f647563745f736869656c6465645f636f6e666967000000000066"
	raw, _ := hex.DecodeString(rawHex)

	p, err := ParsePacket(raw)
	assert.NoError(t, err)

	msg, err := p.DUMLMessage()
	assert.NoError(t, err)
	assert.Equal(t, duml.InterfaceID{Sender: duml.ComponentIDApp, Receiver: 0x28}, msg.Interface)
	assert.Equal(t, uint16(0xde94), uint16(msg.ID))
}

func TestPacket_Video(t *testing.T) {
	// Frame 677 from wlan0.pcap
	// Type: 0x85 (Video)
	rawHex := "c08547a8484202a2384248420000000002050000000001ffaa1b00009011ce00a187fc270000000165b8205bff10307ff7edf59ec81396617513c4e8322d4ac6fab2930da0160c73e258b1cc6c95adfbdcaf73e10904d22487a8e2a19ddf707bc676935029fd9ef08f1e9a7aeec0200c5bf"
	raw, _ := hex.DecodeString(rawHex)

	p, err := ParsePacket(raw)
	assert.NoError(t, err)
	assert.Equal(t, MessageTypeStandard, p.Type)
	assert.Equal(t, WhTypeVideo, p.WhType)
	// Verify NAL start code in payload
	assert.Contains(t, hex.EncodeToString(p.Payload), "0000000165")
}
