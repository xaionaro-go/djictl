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

type CommandID uint8

func (id CommandID) String() string {
	return fmt.Sprintf("0x%02X", uint8(id))
}

type MessageType struct {
	Flags  MessageTypeFlags
	CmdSet CommandSet
	CmdID  CommandID
}

func NewRequest(set CommandSet, id CommandID) MessageType {
	return MessageType{
		Flags:  MessageTypeFlagAckRequired,
		CmdSet: set,
		CmdID:  id,
	}
}

func NewResponse(set CommandSet, id CommandID) MessageType {
	return MessageType{
		Flags:  MessageTypeFlagResponse,
		CmdSet: set,
		CmdID:  id,
	}
}

func (t MessageType) WithFlags(f MessageTypeFlags) MessageType {
	t.Flags = f
	return t
}

const (
	CommandIDGetVersion   CommandID = 0x1E
	CommandIDGetProductID CommandID = 0x0D

	CommandIDOsmoBroadcastConfig CommandID = 0x08
	CommandIDStartStopStreaming  CommandID = 0x8E
	CommandIDPrepareToLiveStream CommandID = 0xE1
	CommandIDConfigureStreaming  CommandID = 0x78

	CommandIDVideoStreamSubscribe   CommandID = 0x3C
	CommandIDVideoStreamUnsubscribe CommandID = 0x3D
	CommandIDGogglesMode            CommandID = 0x3D

	CommandIDRemoteControllerSimulatorData CommandID = 0x24

	CommandIDBatteryStatus  CommandID = 0x02
	CommandIDGetBatteryInfo CommandID = 0x03

	CommandIDFlightStickData CommandID = 0x02
	CommandIDMotorControl    CommandID = 0x21

	CommandIDFCCSupport   CommandID = 0xDE
	CommandIDGetSerialNum CommandID = 0x0A

	CommandIDMaybeStatus    CommandID = 0x05
	CommandIDMaybeKeepAlive CommandID = 0x27

	CommandIDPairingStage2      CommandID = 0x32
	CommandIDPairingStarted     CommandID = 0x80
	CommandIDSetPairingPIN      CommandID = 0x45
	CommandIDPairingPINApproved CommandID = 0x46
	CommandIDConnectToWiFi      CommandID = 0x47
	CommandIDStartScanningWiFi  CommandID = 0xAB
	CommandIDWiFiScanReport     CommandID = 0xAC
	CommandIDCameraAPInfo       CommandID = 0x07
	CommandIDCameraAPInfoResult CommandID = 0x0E
)

var (
	// See: https://github.com/xaionaro/reverse-engineering-dji

	// --- General / Info (Set 0x01) ---
	MessageTypeGetVersion   = NewRequest(CommandSetInfo, CommandIDGetVersion)
	MessageTypeGetProductID = NewRequest(CommandSetInfo, CommandIDGetProductID)

	// --- Video / Camera (Set 0x02) ---
	MessageTypeOsmoBroadcastConfig       = NewRequest(CommandSetCamera, CommandIDOsmoBroadcastConfig)
	MessageTypeStartStopStreaming        = NewRequest(CommandSetCamera, CommandIDStartStopStreaming)
	MessageTypeStartStopStreamingResult  = NewResponse(CommandSetCamera, CommandIDStartStopStreaming)
	MessageTypePrepareToLiveStream       = NewRequest(CommandSetCamera, CommandIDPrepareToLiveStream)
	MessageTypePrepareToLiveStreamResult = NewResponse(CommandSetCamera, CommandIDPrepareToLiveStream).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypeConfigureStreaming        = NewRequest(CommandSetConfig, CommandIDConfigureStreaming)

	// --- Goggles 2 / USB (Set 0x02) ---
	MessageTypeVideoStreamSubscribe   = NewRequest(CommandSetCamera, CommandIDVideoStreamSubscribe)
	MessageTypeVideoStreamUnsubscribe = NewRequest(CommandSetCamera, CommandIDVideoStreamUnsubscribe)
	MessageTypeGogglesMode            = NewRequest(CommandSetFlightController, CommandIDGogglesMode)

	// --- Remote Controller / Simulator (Set 0x06) ---
	MessageTypeRemoteControllerSimulatorData = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetRemoteController, CmdID: CommandIDRemoteControllerSimulatorData}

	// --- Battery / Power (Set 0x0D) ---
	MessageTypeBatteryStatus  = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetBattery, CmdID: CommandIDBatteryStatus}
	MessageTypeGetBatteryInfo = NewRequest(CommandSetBattery, CommandIDGetBatteryInfo)

	// --- Flight Control (Set 0x03) ---
	MessageTypeFlightStickData = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetFlightController, CmdID: CommandIDFlightStickData}
	MessageTypeMotorControl    = NewRequest(CommandSetFlightController, CommandIDMotorControl)

	// --- Common / Config (Set 0x00) ---
	MessageTypeFCCSupport   = NewRequest(CommandSetGeneral, CommandIDFCCSupport)
	MessageTypeGetSerialNum = NewRequest(CommandSetGeneral, CommandIDGetSerialNum)

	// --- InterfaceID: FlightControllerToApp (0x0402) ---
	MessageTypeMaybeStatus    = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: CommandIDMaybeStatus}
	MessageTypeMaybeKeepAlive = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetGimbal, CmdID: CommandIDMaybeKeepAlive}

	// --- InterfaceID: AppToPairer ---
	MessageTypePairingStage2           = NewRequest(CommandSetGeneral, CommandIDPairingStage2)
	MessageTypePairingStarted          = MessageType{Flags: MessageTypeFlagRequest, CmdSet: CommandSetCamera, CmdID: CommandIDPairingStarted}
	MessageTypeSetPairingPIN           = NewRequest(CommandSetWiFi, CommandIDSetPairingPIN)
	MessageTypePairingStatus           = NewResponse(CommandSetWiFi, CommandIDSetPairingPIN).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypePairingPINApproved      = NewRequest(CommandSetWiFi, CommandIDPairingPINApproved)
	MessageTypePairingStage1           = NewResponse(CommandSetWiFi, CommandIDPairingPINApproved).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypeConnectToWiFi           = NewRequest(CommandSetWiFi, CommandIDConnectToWiFi)
	MessageTypeConnectToWiFiResult     = NewResponse(CommandSetWiFi, CommandIDConnectToWiFi).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypeStartScanningWiFi       = NewRequest(CommandSetWiFi, CommandIDStartScanningWiFi)
	MessageTypeStartScanningWiFiResult = NewResponse(CommandSetWiFi, CommandIDStartScanningWiFi).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypeWiFiScanReport          = NewRequest(CommandSetWiFi, CommandIDWiFiScanReport)
	MessageTypeCameraAPInfo            = NewRequest(CommandSetWiFi, CommandIDCameraAPInfo)
	MessageTypeCameraAPInfoResultSSID  = NewResponse(CommandSetWiFi, CommandIDCameraAPInfo).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
	MessageTypeCameraAPInfoResultPSK   = NewResponse(CommandSetWiFi, CommandIDCameraAPInfoResult).WithFlags(MessageTypeFlagResponse | MessageTypeFlagAckRequired)
)

func (t MessageType) GetFlags() uint8 {
	return uint8(t.Flags)
}

func (t MessageType) GetCmdSet() uint8 {
	return uint8(t.CmdSet)
}

func (t MessageType) GetCmdID() uint8 {
	return uint8(t.CmdID)
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
		return fmt.Sprintf("flags:%s set:%s id:%s", t.Flags, t.CmdSet, t.CmdID)
	}
}

func (t *MessageType) ParseFrom(r io.Reader) error {
	var b [3]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return err
	}
	t.Flags = MessageTypeFlags(b[0])
	t.CmdSet = CommandSet(b[1])
	t.CmdID = CommandID(b[2])
	return nil
}

func (t MessageType) Bytes() []byte {
	return []byte{uint8(t.Flags), uint8(t.CmdSet), uint8(t.CmdID)}
}
