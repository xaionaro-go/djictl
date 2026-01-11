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
		_, err := s.RequestStartScanningWiFi(ctx)
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

	msg, err := s.RequestConnectToWiFi(ctx, ssid, psk)
	if err != nil {
		return fmt.Errorf("unable to send the duml.Message: %w", err)
	}

	logger.Debugf(ctx, "received a report about connecting to WiFi: %#+v", msg)
	if !bytes.Equal(msg.Payload, []byte{0, 0}) {
		return fmt.Errorf("unable to connect to WiFi, payload should be 0000, but received %X", msg.Payload)
	}

	return nil
}

func (s *InterfaceAppToWiFiGroundStation) RequestConnectToWiFi(
	ctx context.Context,
	ssid string,
	psk string,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "RequestConnectToWiFi")
	defer func() { logger.Tracef(ctx, "/RequestConnectToWiFi: %v", _err) }()
	msg := s.GetMessageConnectToWiFi(ssid, psk)
	return s.Device().Request(ctx, msg, true)
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

func (s *InterfaceAppToWiFiGroundStation) RequestStartScanningWiFi(
	ctx context.Context,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "RequestStartScanningWiFi")
	defer func() { logger.Tracef(ctx, "/RequestStartScanningWiFi: %v", _err) }()
	msg := s.GetMessageStartScanningWiFi()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessageStartScanningWiFi() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDStartScanningWiFi,
		Type:      duml.MessageTypeStartScanningWiFi,
		Payload:   nil,
	}
}
