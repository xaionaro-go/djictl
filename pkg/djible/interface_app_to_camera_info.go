package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToCamera) GetBatteryInfo(ctx context.Context) (*duml.Message, error) {
	msg := &duml.Message{
		Interface: s.InterfaceID(),
		Type:      duml.MessageTypeGetBatteryInfo,
	}
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToCamera) GetVersion(ctx context.Context) (*duml.Message, error) {
	msg := &duml.Message{
		Interface: s.InterfaceID(),
		Type:      duml.MessageTypeGetVersion,
	}
	return s.Device().Request(ctx, msg, true)
}
