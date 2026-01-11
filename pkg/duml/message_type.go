package duml

import (
	"fmt"
	"io"
	"strings"
)

type MessageTypeFlags uint8

const (
	MessageTypeFlagRequest     MessageTypeFlags = 0x00
	MessageTypeFlagAckRequired MessageTypeFlags = 0x40
	MessageTypeFlagResponse    MessageTypeFlags = 0x80
)

func (f MessageTypeFlags) String() string {
	var parts []string
	if f&MessageTypeFlagResponse != 0 {
		parts = append(parts, "Response")
	} else {
		parts = append(parts, "Request")
	}
	if f&MessageTypeFlagAckRequired != 0 {
		if f&MessageTypeFlagResponse != 0 {
			parts = append(parts, "Ack")
		} else {
			parts = append(parts, "AckRequired")
		}
	}
	return strings.Join(parts, "|")
}

type CommandSet uint8

const (
	CommandSetGeneral          CommandSet = 0x00
	CommandSetInfo             CommandSet = 0x01
	CommandSetCamera           CommandSet = 0x02
	CommandSetFlightController CommandSet = 0x03
	CommandSetGimbal           CommandSet = 0x04
	CommandSetRemoteController CommandSet = 0x06
	CommandSetWiFi             CommandSet = 0x07
	CommandSetConfig           CommandSet = 0x08
	CommandSetVision           CommandSet = 0x0a
	CommandSetBattery          CommandSet = 0x0d
)

func (s CommandSet) String() string {
	switch s {
	case CommandSetGeneral:
		return "General"
	case CommandSetInfo:
		return "Info"
	case CommandSetCamera:
		return "Camera"
	case CommandSetFlightController:
		return "FlightController"
	case CommandSetGimbal:
		return "Gimbal"
	case CommandSetRemoteController:
		return "RemoteController"
	case CommandSetWiFi:
		return "WiFi"
	case CommandSetConfig:
		return "Config"
	case CommandSetVision:
		return "Vision"
	case CommandSetBattery:
		return "Battery"
	default:
		return fmt.Sprintf("0x%02X", uint8(s))
	}
}

type MessageType struct {
	Flags  MessageTypeFlags
	CmdSet CommandSet
	CmdID  uint8
}

var (
	// See: https://github.com/xaionaro/reverse-engineering-dji

	// --- General / Info (Set 0x01) ---
	MessageTypeGetVersion   = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetInfo, CmdID: 0x1E}
	MessageTypeGetProductID = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetInfo, CmdID: 0x0D}

	// --- Video / Camera (Set 0x02) ---
	MessageTypeOsmoBroadcastConfig       = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0x08}
	MessageTypeStartStopStreaming        = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0x8E}
	MessageTypeStartStopStreamingResult  = MessageType{Flags: MessageTypeFlagResponse, CmdSet: CommandSetCamera, CmdID: 0x8E}
	MessageTypePrepareToLiveStream       = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0xE1}
	MessageTypePrepareToLiveStreamResult = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0xE1}
	MessageTypeConfigureStreaming        = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetConfig, CmdID: 0x78}

	// --- Goggles 2 / USB (Set 0x02) ---
	MessageTypeVideoStreamSubscribe   = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0x3C}
	MessageTypeVideoStreamUnsubscribe = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetCamera, CmdID: 0x3D}
	MessageTypeGogglesMode            = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetFlightController, CmdID: 0x3D}

	// --- Remote Controller / Simulator (Set 0x06) ---
	MessageTypeRemoteControllerSimulatorData = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetRemoteController, CmdID: 0x24}

	// --- Battery / Power (Set 0x0D) ---
	MessageTypeBatteryStatus  = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetBattery, CmdID: 0x02}
	MessageTypeGetBatteryInfo = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetBattery, CmdID: 0x03}

	// --- Flight Control (Set 0x03) ---
	MessageTypeFlightStickData = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetFlightController, CmdID: 0x02}
	MessageTypeMotorControl    = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetFlightController, CmdID: 0x21}

	// --- Common / Config (Set 0x00) ---
	MessageTypeFCCSupport   = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetGeneral, CmdID: 0xDE}
	MessageTypeGetSerialNum = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetGeneral, CmdID: 0x0A}

	// --- InterfaceID: FlightControllerToApp (0x0402) ---
	MessageTypeMaybeStatus    = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: 0x05}
	MessageTypeMaybeKeepAlive = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: 0x27}

	// --- InterfaceID: AppToPairer ---
	MessageTypePairingStage2           = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetGeneral, CmdID: 0x32}
	MessageTypePairingStarted          = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetCamera, CmdID: 0x80}
	MessageTypeSetPairingPIN           = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x45}
	MessageTypePairingStatus           = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x45}
	MessageTypePairingPINApproved      = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x46}
	MessageTypePairingStage1           = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x46}
	MessageTypeConnectToWiFi           = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x47}
	MessageTypeConnectToWiFiResult     = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x47}
	MessageTypeStartScanningWiFi       = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0xAB}
	MessageTypeStartScanningWiFiResult = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0xAB}
	MessageTypeWiFiScanReport          = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0xAC}
	MessageTypeCameraAPInfo            = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x07}
	MessageTypeCameraAPInfoResultSSID  = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x07}
	MessageTypeCameraAPInfoResultPSK   = MessageType{Flags: MessageTypeFlagResponse | MessageTypeFlagAckRequired, CmdSet: CommandSetWiFi, CmdID: 0x0E}

	MessageTypeUnknown0 = MessageType{Flags: MessageTypeFlagAckRequired, CmdSet: CommandSetGeneral, CmdID: 0x81}
	MessageTypeUnknown1 = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGeneral, CmdID: 0xF1}
	MessageTypeUnknown2 = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetCamera, CmdID: 0xDC}
	MessageTypeUnknown3 = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: 0x1C}
	MessageTypeUnknown4 = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: 0x38}
	MessageTypeUnknown5 = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetWiFi, CmdID: 0x45}
)

func (t MessageType) GetFlags() uint8 {
	return uint8(t.Flags)
}

func (t MessageType) GetCmdSet() uint8 {
	return uint8(t.CmdSet)
}

func (t MessageType) GetCmdID() uint8 {
	return t.CmdID
}

func (t MessageType) String() string {
	switch t {
	case MessageTypeGetVersion:
		return "get_version"
	case MessageTypeGetProductID:
		return "get_product_id"
	case MessageTypeOsmoBroadcastConfig:
		return "osmo_broadcast_config"
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
	case MessageTypeStartStopStreaming:
		return "start_OR_stop_streaming"
	case MessageTypeStartStopStreamingResult:
		return "start_OR_stop_streaming_result"
	case MessageTypeBatteryStatus:
		return "battery_status"
	case MessageTypeWiFiScanReport:
		return "wifi_scan_results"
	case MessageTypeStartScanningWiFi:
		return "start_scanning_wifi"
	case MessageTypeStartScanningWiFiResult:
		return "start_scanning_wifi_result"
	case MessageTypeCameraAPInfo:
		return "camera_ap_info"
	case MessageTypeCameraAPInfoResultSSID:
		return "camera_ap_info_result_ssid"
	case MessageTypeCameraAPInfoResultPSK:
		return "camera_ap_info_result_psk"
	case MessageTypeRemoteControllerSimulatorData:
		return "remote_controller_simulator_data"
	default:
		return fmt.Sprintf("flags:%s set:%s id:%02X", t.Flags, t.CmdSet, t.CmdID)
	}
}

func (t *MessageType) ParseFrom(r io.Reader) error {
	var b [3]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return err
	}
	t.Flags = MessageTypeFlags(b[0])
	t.CmdSet = CommandSet(b[1])
	t.CmdID = b[2]
	return nil
}

func (t MessageType) Bytes() []byte {
	return []byte{uint8(t.Flags), uint8(t.CmdSet), t.CmdID}
}
