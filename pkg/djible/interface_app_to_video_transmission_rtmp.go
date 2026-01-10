package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToVideoTransmission) ConfigureRTMP(ctx context.Context, url string, enable bool) error {
	msg := duml.NewBroadcastMessage(enable, url)
	return s.Device().SendMessage(ctx, msg, true)
}
