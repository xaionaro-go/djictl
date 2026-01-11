package djible

import (
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

type InterfaceAppToBattery Device

func (d *Device) AppToBattery() *InterfaceAppToBattery {
	return (*InterfaceAppToBattery)(d)
}

func (s *InterfaceAppToBattery) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToBattery
}

func (s *InterfaceAppToBattery) Device() *Device {
	return (*Device)(s)
}

func (s *InterfaceAppToBattery) GetInfo(ctx context.Context) (*duml.BatteryStatus, error) {
	msg := &duml.Message{
		Interface: s.InterfaceID(),
		Type:      duml.MessageTypeGetBatteryInfo,
	}
	_, err := s.Device().Request(ctx, msg, true)
	if err != nil {
		return nil, fmt.Errorf("unable to send GetBatteryInfo message: %w", err)
	}
	logger.Debugf(ctx, "waiting for the battery status")

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypeBatteryStatus):
		logger.Debugf(ctx, "received a report about battery info: %#+v", msg)
		status, err := duml.ParseBatteryStatus(ctx, msg.Payload)
		if err != nil {
			return nil, fmt.Errorf("unable to parse battery status: %w", err)
		}
		return status, nil
	}
}
