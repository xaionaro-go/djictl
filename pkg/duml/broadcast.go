package duml

import "bytes"

type BroadcastPlatform uint8

const (
	BroadcastPlatformRTMP = BroadcastPlatform(2)
)

type BroadcastConfig struct {
	Enabled  bool
	Platform BroadcastPlatform
	URL      string
}

func (c *BroadcastConfig) Payload() []byte {
	var buf bytes.Buffer
	if c.Enabled {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	buf.WriteByte(uint8(c.Platform))

	buf.Write(PackURL(c.URL))
	return buf.Bytes()
}

func NewBroadcastMessage(enabled bool, url string) *Message {
	config := &BroadcastConfig{
		Enabled:  enabled,
		Platform: BroadcastPlatformRTMP,
		URL:      url,
	}
	return &Message{
		Interface: InterfaceIDAppToCamera,
		ID:        MessageIDStartStreaming,
		Type:      MessageTypeOsmoBroadcastConfig,
		Payload:   config.Payload(),
	}
}
