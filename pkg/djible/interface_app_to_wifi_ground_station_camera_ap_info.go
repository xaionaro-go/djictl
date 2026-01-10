package djible

import (
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToWiFiGroundStation) CameraAPInfo(
	ctx context.Context,
) (string, string, error) {
	logger.Tracef(ctx, "CameraAPInfo")
	defer func() { logger.Tracef(ctx, "/CameraAPInfo") }()

	err := s.Device().SendMessage(ctx, s.GetMessageCameraAPInfo(), true)
	if err != nil {
		return "", "", fmt.Errorf("unable to send the duml.Message: %w", err)
	}

	var ssid, psk string
	for ssid == "" || psk == "" {
		select {
		case <-ctx.Done():
			return "", "", ctx.Err()
		case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypeCameraAPInfoResultSSID):
			logger.Debugf(ctx, "received SSID: %X", msg.Payload)
			var err error
			ssid, err = duml.UnpackStringU16BE(msg.Payload)
			if err != nil {
				return "", "", fmt.Errorf("unable to unpack SSID: %w", err)
			}
		case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypeCameraAPInfoResultPSK):
			logger.Debugf(ctx, "received PSK: %X", msg.Payload)
			var err error
			psk, err = duml.UnpackStringU16BE(msg.Payload)
			if err != nil {
				return "", "", fmt.Errorf("unable to unpack PSK: %w", err)
			}
		}
	}

	return ssid, psk, nil
}

func (s *InterfaceAppToWiFiGroundStation) GetMessageCameraAPInfo() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDCameraAPInfo,
		Type:      duml.MessageTypeCameraAPInfo,
		Payload:   []byte{0x20},
	}
}
