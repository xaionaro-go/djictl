package djictl

import "fmt"

type MessageID uint16

const (
	MessageIDPairingStarted            = MessageID(0x7911)
	MessageIDSetPairingPIN             = MessageID(0x72AA)
	MessageIDPairingStage1             = MessageID(0x0400)
	MessageIDPairingStage2             = MessageID(0x74AA)
	MessageIDPrepareToLiveStreamStage1 = MessageID(0xFEAB)
	MessageIDPrepareToLiveStreamStage2 = MessageID(0xFFAB)
	MessageIDStartScanningWiFi         = MessageID(0x8EBB)
	MessageIDConnectToWifi             = MessageID(0x98BB)
	MessageIDConfigureStreaming        = MessageID(0xB3BB)
	MessageIDStartStreaming            = MessageID(0xB4BB)
	MessageIDStopStreaming             = MessageID(0xB5BB)
)

func (id MessageID) String() string {
	switch id {
	case MessageIDSetPairingPIN:
		return "pair"
	case MessageIDPrepareToLiveStreamStage1:
		return "prepare_to_live_stream_stage1"
	case MessageIDPrepareToLiveStreamStage2:
		return "prepare_to_live_stream_stage2"
	case MessageIDConnectToWifi:
		return "connect_to_wifi"
	case MessageIDConfigureStreaming:
		return "configure_streaming"
	case MessageIDStartStreaming:
		return "start_streaming"
	case MessageIDStopStreaming:
		return "stop_streaming"
	default:
		return fmt.Sprintf("%04X", uint16(id))
	}
}
