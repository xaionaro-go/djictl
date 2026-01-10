package duml

import (
	"fmt"
)

type InterfaceID struct {
	Sender   ComponentID
	Receiver ComponentID
}

type ComponentID uint8

const (
	ComponentIDCamera            ComponentID = 0x01
	ComponentIDApp               ComponentID = 0x02
	ComponentIDGimbal            ComponentID = 0x03
	ComponentIDFlightController  ComponentID = 0x04
	ComponentIDWiFiAir           ComponentID = 0x05
	ComponentIDRemoteController  ComponentID = 0x06
	ComponentIDWiFiGroundStation ComponentID = 0x07
	ComponentIDVideoTransmission ComponentID = 0x08
	ComponentIDBattery           ComponentID = 0x09
	ComponentIDGimbal2           ComponentID = 0x0e
	ComponentIDVision            ComponentID = 0x11
	ComponentIDGoggles           ComponentID = 0x17
	ComponentIDPairer            ComponentID = 0x88
)

func (id ComponentID) String() string {
	switch id {
	case ComponentIDCamera:
		return "camera"
	case ComponentIDApp:
		return "app"
	case ComponentIDGimbal:
		return "gimbal"
	case ComponentIDFlightController:
		return "flight_controller"
	case ComponentIDRemoteController:
		return "remote_controller"
	case ComponentIDBattery:
		return "battery"
	case ComponentIDGoggles:
		return "goggles"
	case ComponentIDWiFiGroundStation:
		return "wifi_ground_station"
	case ComponentIDVideoTransmission:
		return "video_transmission"
	case ComponentIDPairer:
		return "pairer"
	default:
		return fmt.Sprintf("%02X", uint8(id))
	}
}

var (
	InterfaceIDAppToApp               = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDApp}
	InterfaceIDAppToCamera            = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDCamera}
	InterfaceIDAppToGimbal            = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDGimbal}
	InterfaceIDAppToFlightController  = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDFlightController}
	InterfaceIDAppToRemoteController  = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDRemoteController}
	InterfaceIDAppToWiFiGroundStation = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDWiFiGroundStation}
	InterfaceIDAppToVideoTransmission = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDVideoTransmission}
	InterfaceIDAppToBattery           = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDBattery}
	InterfaceIDAppToGoggles           = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDGoggles}
	InterfaceIDAppToPairer            = InterfaceID{Sender: ComponentIDApp, Receiver: ComponentIDPairer}
	InterfaceIDFlightControllerToApp  = InterfaceID{Sender: ComponentIDFlightController, Receiver: ComponentIDApp}
)

func (id InterfaceID) String() string {
	return fmt.Sprintf("%s->%s", id.Sender, id.Receiver)
}
