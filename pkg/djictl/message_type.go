package djictl

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MessageType uint32

const (
	// See: https://github.com/xaionaro/reverse-engineering-dji

	// SubsystemID: Configurer
	MessageTypeConfigure = MessageType(0x40028E)

	// SubsystemID: PrePairer (0x0402)
	MessageTypeMaybeStatus    = MessageType(0x000405)
	MessageTypeMaybeKeepAlive = MessageType(0x000427)

	// SubsystemID: Pairer
	MessageTypePairingStage2           = MessageType(0x400032)
	MessageTypePairingStarted          = MessageType(0x000280)
	MessageTypeSetPairingPIN           = MessageType(0x400745)
	MessageTypePairingStatus           = MessageType(0xC00745)
	MessageTypePairingPINApproved      = MessageType(0x400746)
	MessageTypePairingStage1           = MessageType(0xC00746)
	MessageTypeConnectToWiFi           = MessageType(0x400747)
	MessageTypeConnectToWiFiResult     = MessageType(0xC00747)
	MessageTypeStartScanningWiFi       = MessageType(0x4007AB)
	MessageTypeStartScanningWiFiResult = MessageType(0xC007AB)
	MessageTypeWiFiScanReport          = MessageType(0x4007AC)

	// SubsystemID: Streamer
	MessageTypeStartStopStreaming        = MessageType(0x40028E)
	MessageTypeStartStopStreamingResult  = MessageType(0x80028E)
	MessageTypePrepareToLiveStream       = MessageType(0x4002E1)
	MessageTypePrepareToLiveStreamResult = MessageType(0xC002E1)
	MessageTypeConfigureStreaming        = MessageType(0x400878)
	MessageTypeStreamingStatus           = MessageType(0x000D02)

	MessageTypeUnknown0 = MessageType(0x400081)
	MessageTypeUnknown1 = MessageType(0x0000F1)
	MessageTypeUnknown2 = MessageType(0x0002DC)
	MessageTypeUnknown3 = MessageType(0x00041C)
	MessageTypeUnknown4 = MessageType(0x000438)
	MessageTypeUnknown5 = MessageType(0x000745)
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeMaybeStatus:
		return "status?"
	case MessageTypeMaybeKeepAlive:
		return "keep_alive?"
	case MessageTypePairingStarted:
		return "pairing_started"
	case MessageTypeSetPairingPIN:
		return "set_pairing_pin"
	case MessageTypePairingStatus:
		return "pairing_status"
	case MessageTypePairingPINApproved:
		return "pairing_pin_approved"
	case MessageTypePairingStage1:
		return "pairing_stage1"
	case MessageTypePairingStage2:
		return "pairing_stage2"
	case MessageTypeConnectToWiFi:
		return "connect_to_wifi"
	case MessageTypePrepareToLiveStream:
		return "prepare_to_live_stream"
	case MessageTypePrepareToLiveStreamResult:
		return "prepare_to_live_stream_result"
	case MessageTypeConfigureStreaming:
		return "configure_stream"
	case MessageTypeStreamingStatus:
		return "streaming_status"
	case MessageTypeStartStopStreaming:
		return "start_OR_stop_streaming"
	case MessageTypeStartStopStreamingResult:
		return "start_OR_stop_streaming_result"
	case MessageTypeWiFiScanReport:
		return "wifi_scan_results"
	case MessageTypeStartScanningWiFi:
		return "start_scanning_wifi"
	case MessageTypeStartScanningWiFiResult:
		return "start_scanning_wifi_result"
	default:
		return fmt.Sprintf("%04X", uint16(t))
	}
}

func (t *MessageType) ParseFrom(r io.Reader) error {
	var b [4]byte
	n, err := r.Read(b[1:])
	if err != nil {
		return err
	}
	if n != 3 {
		return fmt.Errorf("%w: expected 3, but read %d", io.ErrShortBuffer, n)
	}
	v := binary.BigEndian.Uint32(b[:])
	*t = MessageType(v)
	return nil
}

func (t MessageType) Bytes() []byte {
	var r [4]byte
	binary.BigEndian.PutUint32(r[:], uint32(t))
	return r[1:]
}
