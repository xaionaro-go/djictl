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
	CommandSetCore             CommandSet = 0x00
	CommandSetInfo             CommandSet = 0x01
	CommandSetCamera           CommandSet = 0x02
	CommandSetFlightController CommandSet = 0x03
	CommandSetUnknown0         CommandSet = 0x04
	CommandSetRemoteController CommandSet = 0x06
	CommandSetWiFi             CommandSet = 0x07
	CommandSetConfig           CommandSet = 0x08
	CommandSetVision           CommandSet = 0x0a
	CommandSetBattery          CommandSet = 0x0d
)

func (s CommandSet) String() string {
	switch s {
	case CommandSetCore:
		return "Core"
	case CommandSetInfo:
		return "Info"
	case CommandSetCamera:
		return "Camera"
	case CommandSetFlightController:
		return "FlightController"
	case CommandSetUnknown0:
		return "Unknown0"
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

func MessageTypeNotification(set CommandSet, id CommandID, flags ...MessageTypeFlags) MessageType {
	var combinedFlags MessageTypeFlags
	for _, f := range flags {
		combinedFlags |= f
	}
	return MessageType{
		Flags:  combinedFlags,
		CmdSet: set,
		CmdID:  id,
	}
}

func MessageTypeRequest(set CommandSet, id CommandID, flags ...MessageTypeFlags) MessageType {
	var combinedFlags MessageTypeFlags = MessageTypeFlagAckRequired
	for _, f := range flags {
		combinedFlags |= f
	}
	return MessageType{
		Flags:  combinedFlags,
		CmdSet: set,
		CmdID:  id,
	}
}

func MessageTypeResponse(set CommandSet, id CommandID, flags ...MessageTypeFlags) MessageType {
	var combinedFlags MessageTypeFlags = MessageTypeFlagResponse
	for _, f := range flags {
		combinedFlags |= f
	}
	return MessageType{
		Flags:  combinedFlags,
		CmdSet: set,
		CmdID:  id,
	}
}

func (t MessageType) WithFlags(f MessageTypeFlags) MessageType {
	t.Flags = f
	return t
}

const (
	// --- Core (Set 0x00) ---
	CommandIDGetSerialNum  CommandID = 0x0A
	CommandIDHeartbeat     CommandID = 0x2B
	CommandIDPairingStage2 CommandID = 0x32
	CommandIDParameterPush CommandID = 0x99
	CommandIDFCCSupport    CommandID = 0xDE

	// --- Info (Set 0x01) ---
	CommandIDGetProductID CommandID = 0x0D
	CommandIDGetVersion   CommandID = 0x1E

	// --- Video / Camera (Set 0x02) ---
	CommandIDTakeRecord             CommandID = 0x02
	CommandIDGogglesModeToggle      CommandID = 0x06
	CommandIDOsmoBroadcastConfig    CommandID = 0x08
	CommandIDVideoStreamSubscribe   CommandID = 0x3C
	CommandIDVideoStreamUnsubscribe CommandID = 0x3D
	CommandIDPairingStarted         CommandID = 0x80
	CommandIDStartStopStreaming     CommandID = 0x8E
	CommandIDPrepareToLiveStream    CommandID = 0xE1

	// --- Flight Control (Set 0x03) ---
	CommandIDFlightStickData CommandID = 0x02
	CommandIDMotorControl    CommandID = 0x21
	CommandIDGogglesMode     CommandID = 0x3D

	// --- Gimbal (Set 0x04) ---
	CommandIDMaybeStatus CommandID = 0x05
	CommandIDKeepAlive   CommandID = 0x27

	// --- Remote Controller (Set 0x06) ---
	CommandIDRemoteControllerSimulatorData CommandID = 0x24

	// --- WiFi (Set 0x07) ---
	CommandIDCameraAPInfo       CommandID = 0x07
	CommandIDCameraAPPSK        CommandID = 0x0E
	CommandIDSetPairingPIN      CommandID = 0x45
	CommandIDPairingPINApproved CommandID = 0x46
	CommandIDConnectToWiFi      CommandID = 0x47
	CommandIDStartScanningWiFi  CommandID = 0xAB
	CommandIDWiFiScanReport     CommandID = 0xAC

	// --- Config (Set 0x08) ---
	CommandIDConfigureStreaming CommandID = 0x78

	// --- Battery (Set 0x0D) ---
	CommandIDBatteryStatus  CommandID = 0x02
	CommandIDGetBatteryInfo CommandID = 0x03
)

var (
	// See: https://github.com/xaionaro/reverse-engineering-dji

	// --- Common / Config (Set 0x00) ---
	MessageTypeGetSerialNum  = MessageTypeRequest(CommandSetCore, CommandIDGetSerialNum)
	MessageTypeHeartbeat     = MessageTypeRequest(CommandSetCore, CommandIDHeartbeat)
	MessageTypePairingStage2 = MessageTypeRequest(CommandSetCore, CommandIDPairingStage2)
	MessageTypeParameterPush = MessageTypeNotification(CommandSetCore, CommandIDParameterPush)
	MessageTypeFCCSupport    = MessageTypeRequest(CommandSetCore, CommandIDFCCSupport)

	// --- General / Info (Set 0x01) ---
	MessageTypeGetProductID = MessageTypeRequest(CommandSetInfo, CommandIDGetProductID)
	MessageTypeGetVersion   = MessageTypeRequest(CommandSetInfo, CommandIDGetVersion)

	// --- Video / Camera (Set 0x02) (also suspected to be reusable for Goggles 2 / USB) ---
	MessageTypeGogglesModeToggle         = MessageTypeRequest(CommandSetCamera, CommandIDGogglesModeToggle)
	MessageTypeOsmoBroadcastConfig       = MessageTypeRequest(CommandSetCamera, CommandIDOsmoBroadcastConfig)
	MessageTypeVideoStreamSubscribe      = MessageTypeRequest(CommandSetCamera, CommandIDVideoStreamSubscribe)
	MessageTypeVideoStreamUnsubscribe    = MessageTypeRequest(CommandSetCamera, CommandIDVideoStreamUnsubscribe)
	MessageTypePairingStarted            = MessageTypeNotification(CommandSetCamera, CommandIDPairingStarted)
	MessageTypeStartStopStreaming        = MessageTypeRequest(CommandSetCamera, CommandIDStartStopStreaming)
	MessageTypeStartStopStreamingResult  = MessageTypeResponse(CommandSetCamera, CommandIDStartStopStreaming)
	MessageTypePrepareToLiveStream       = MessageTypeRequest(CommandSetCamera, CommandIDPrepareToLiveStream)
	MessageTypePrepareToLiveStreamResult = MessageTypeResponse(CommandSetCamera, CommandIDPrepareToLiveStream, MessageTypeFlagAckRequired)

	// --- Flight Control (Set 0x03) ---
	MessageTypeFlightStickData = MessageTypeNotification(CommandSetFlightController, CommandIDFlightStickData)
	MessageTypeMotorControl    = MessageTypeRequest(CommandSetFlightController, CommandIDMotorControl)
	MessageTypeGogglesMode     = MessageTypeRequest(CommandSetFlightController, CommandIDGogglesMode)

	// --- InterfaceID: FlightControllerToApp (0x0402) ---
	MessageTypeUnknown0MaybeStatus = MessageTypeNotification(CommandSetUnknown0, CommandIDMaybeStatus)
	MessageTypeKeepAlive           = MessageTypeNotification(CommandSetUnknown0, CommandIDKeepAlive)

	// --- Remote Controller / Simulator (Set 0x06) ---
	MessageTypeRemoteControllerSimulatorData = MessageTypeNotification(CommandSetRemoteController, CommandIDRemoteControllerSimulatorData)

	// --- InterfaceID: AppToPairer ---
	MessageTypeCameraAPInfo            = MessageTypeRequest(CommandSetWiFi, CommandIDCameraAPInfo)
	MessageTypeCameraAPInfoResultSSID  = MessageTypeResponse(CommandSetWiFi, CommandIDCameraAPInfo, MessageTypeFlagAckRequired)
	MessageTypeGetCameraAPPSK          = MessageTypeRequest(CommandSetWiFi, CommandIDCameraAPPSK)
	MessageTypeCameraAPInfoResultPSK   = MessageTypeResponse(CommandSetWiFi, CommandIDCameraAPPSK, MessageTypeFlagAckRequired)
	MessageTypeSetPairingPIN           = MessageTypeRequest(CommandSetWiFi, CommandIDSetPairingPIN)
	MessageTypePairingStatus           = MessageTypeResponse(CommandSetWiFi, CommandIDSetPairingPIN, MessageTypeFlagAckRequired)
	MessageTypePairingPINApproved      = MessageTypeRequest(CommandSetWiFi, CommandIDPairingPINApproved)
	MessageTypePairingStage1           = MessageTypeResponse(CommandSetWiFi, CommandIDPairingPINApproved, MessageTypeFlagAckRequired)
	MessageTypeConnectToWiFi           = MessageTypeRequest(CommandSetWiFi, CommandIDConnectToWiFi)
	MessageTypeConnectToWiFiResult     = MessageTypeResponse(CommandSetWiFi, CommandIDConnectToWiFi, MessageTypeFlagAckRequired)
	MessageTypeStartScanningWiFi       = MessageTypeRequest(CommandSetWiFi, CommandIDStartScanningWiFi)
	MessageTypeStartScanningWiFiResult = MessageTypeResponse(CommandSetWiFi, CommandIDStartScanningWiFi, MessageTypeFlagAckRequired)
	MessageTypeWiFiScanReport          = MessageTypeRequest(CommandSetWiFi, CommandIDWiFiScanReport)

	// --- Config (Set 0x08) ---
	MessageTypeConfigureStreaming       = MessageTypeRequest(CommandSetConfig, CommandIDConfigureStreaming)
	MessageTypeConfigureStreamingResult = MessageTypeResponse(CommandSetConfig, CommandIDConfigureStreaming, MessageTypeFlagAckRequired)

	// --- Battery / Power (Set 0x0D) ---
	MessageTypeBatteryStatus  = MessageTypeNotification(CommandSetBattery, CommandIDBatteryStatus)
	MessageTypeGetBatteryInfo = MessageTypeRequest(CommandSetBattery, CommandIDGetBatteryInfo)
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
	case MessageTypeVideoStreamSubscribe:
		return "video_stream_subscribe"
	case MessageTypeVideoStreamUnsubscribe:
		return "video_stream_unsubscribe"
	case MessageTypeGogglesModeToggle:
		return "goggles_mode_toggle"
	case MessageTypeGogglesMode:
		return "goggles_mode"
	case MessageTypeRemoteControllerSimulatorData:
		return "remote_controller_simulator_data"
	case MessageTypeBatteryStatus:
		return "battery_status"
	case MessageTypeGetBatteryInfo:
		return "get_battery_info"
	case MessageTypeFlightStickData:
		return "flight_stick_data"
	case MessageTypeMotorControl:
		return "motor_control"
	case MessageTypeFCCSupport:
		return "fcc_support"
	case MessageTypeGetSerialNum:
		return "get_serial_num"
	case MessageTypeHeartbeat:
		return "heartbeat"
	case MessageTypeParameterPush:
		return "parameter_push"
	case MessageTypeUnknown0MaybeStatus:
		return "gimbal_status"
	case MessageTypeKeepAlive:
		return "keep_alive"
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
	case MessageTypeConfigureStreamingResult:
		return "configure_stream_result"
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
	case MessageTypeCameraAPInfo:
		return "camera_ap_info"
	case MessageTypeGetCameraAPPSK:
		return "get_camera_ap_psk"
	case MessageTypeCameraAPInfoResultSSID:
		return "camera_ap_info_result_ssid"
	case MessageTypeCameraAPInfoResultPSK:
		return "camera_ap_info_result_psk"
	case MessageTypeOsmoBroadcastConfig:
		return "osmo_broadcast_config"
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
