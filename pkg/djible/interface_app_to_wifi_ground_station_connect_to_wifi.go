package djible

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

const (
	wifiWaitForScanReport = false
)

func (s *InterfaceAppToWiFiGroundStation) ConnectToWiFi(
	ctx context.Context,
	ssid string,
	psk string,
) (_err error) {
	logger.Tracef(ctx, "ConnectToWiFi")
	defer func() { logger.Tracef(ctx, "/ConnectToWiFi: %v", _err) }()

	if wifiWaitForScanReport {
		err := s.SendMessageStartScanningWiFi(ctx)
		if err != nil {
			return fmt.Errorf("unable to send the duml.Message: %w", err)
		}

		logger.Debugf(ctx, "waiting WiFi scan results")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypeWiFiScanReport):
			logger.Debugf(ctx, "received a wifi scan result: %#+v", msg)
		}
	}

	err := s.SendMessageConnectToWiFi(ctx, ssid, psk)
	if err != nil {
		return fmt.Errorf("unable to send the duml.Message: %w", err)
	}

	logger.Debugf(ctx, "waiting for connecting to WiFi")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypeConnectToWiFiResult):
		logger.Debugf(ctx, "received a report about connecting to WiFi: %#+v", msg)
		if !bytes.Equal(msg.Payload, []byte{0, 0}) {
			return fmt.Errorf("unable to connect to WiFi, payload should be 0000, but received %X", msg.Payload)
		}
	}

	return nil
}

func (s *InterfaceAppToWiFiGroundStation) SendMessageConnectToWiFi(
	ctx context.Context,
	ssid string,
	psk string,
) (_err error) {
	logger.Tracef(ctx, "SendMessageConnectToWiFi")
	defer func() { logger.Tracef(ctx, "/SendMessageConnectToWiFi: %v", _err) }()
	msg := s.GetMessageConnectToWiFi(ssid, psk)
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessageConnectToWiFi(
	ssid string,
	psk string,
) *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDConnectToWifi,
		Type:      duml.MessageTypeConnectToWiFi,
		Payload: s.GetMessagePayloadConnectToWiFi(
			ssid,
			psk,
		),
	}
}

func (s *InterfaceAppToWiFiGroundStation) GetMessagePayloadConnectToWiFi(
	ssid string,
	psk string,
) []byte {
	var buf bytes.Buffer
	must(buf.Write(duml.PackString(ssid)))
	must(buf.Write(duml.PackString(psk)))
	return buf.Bytes()
}

func (s *InterfaceAppToWiFiGroundStation) SendMessageStartScanningWiFi(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendMessageStartScanningWiFi")
	defer func() { logger.Tracef(ctx, "/SendMessageStartScanningWiFi: %v", _err) }()
	msg := s.GetMessageStartScanningWiFi()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessageStartScanningWiFi() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDStartScanningWiFi,
		Type:      duml.MessageTypeStartScanningWiFi,
		Payload:   nil,
	}
}
