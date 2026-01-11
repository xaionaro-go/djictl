package djible

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

const (
	defaultPINCode = "5160"
)

func (s *InterfaceAppToWiFiGroundStation) Pair(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "Pair")
	defer func() { logger.Tracef(ctx, "/Pair: %v", _err) }()
	err := s.SendRequestStartPairing(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the request to start pairing: %w", err)
	}
	msg, err := s.RequestSetPairingPIN(ctx, defaultPINCode)
	if err != nil {
		return fmt.Errorf("unable to send the message to set the PIN: %w", err)
	}

	logger.Debugf(ctx, "received the pairing info: %#+v", msg)
	if len(msg.Payload) < 2 {
		logger.Errorf(ctx, "the payload size of the pairing status is too small: %d", len(msg.Payload))
	} else {
		if msg.Payload[1] == 0x01 {
			logger.Debugf(ctx, "is already paired")
			return nil
		}
	}

	logger.Debugf(ctx, "waiting for PIN approve")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, duml.MessageTypePairingPINApproved):
		logger.Debugf(ctx, "PIN was approved: %#+v", msg)
	}

	_, err = s.RequestPairingStage1(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the command to pair (stage1): %w", err)
	}
	_, err = s.RequestPairingStage2(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the command to pair (stage2): %w", err)
	}
	return nil
}

func (s *InterfaceAppToWiFiGroundStation) SendRequestStartPairing(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendRequestStartPairing")
	defer func() { logger.Tracef(ctx, "/SendRequestStartPairing: %v", _err) }()
	return s.Device().SendPairingRequest(ctx)
}

func (s *InterfaceAppToWiFiGroundStation) RequestSetPairingPIN(
	ctx context.Context,
	pinCode string,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "RequestSetPairingPIN")
	defer func() { logger.Tracef(ctx, "/RequestSetPairingPIN: %v", _err) }()
	msg := s.GetMessageSetPairingPIN(pinCode)
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessageSetPairingPIN(
	pinCode string,
) *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDSetPairingPIN,
		Type:      duml.MessageTypeSetPairingPIN,
		Payload: s.GetMessagePayloadSetPairingPIN(
			pinCode,
		),
	}
}

func (s *InterfaceAppToWiFiGroundStation) GetMessagePayloadSetPairingPIN(
	pinCode string,
) []byte {
	var buf bytes.Buffer
	must(buf.Write(duml.PackString("001749319286102")))
	must(buf.Write(duml.PackString(pinCode)))
	return buf.Bytes()
}

func (s *InterfaceAppToWiFiGroundStation) ReceiveMessageSetPairingPINResult(
	ctx context.Context,
) (*duml.Message, error) {
	return s.Device().ReceiveMessage(ctx, duml.MessageTypeSetPairingPIN)
}

func (s *InterfaceAppToWiFiGroundStation) RequestPairingStage1(
	ctx context.Context,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "RequestPairingStage1")
	defer func() { logger.Tracef(ctx, "/RequestPairingStage1: %v", _err) }()
	msg := s.GetMessagePairingStage1()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessagePairingStage1() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDPairingStage1,
		Type:      duml.MessageTypePairingStage1,
		Payload:   []byte{0x00},
	}
}

func (s *InterfaceAppToWiFiGroundStation) RequestPairingStage2(
	ctx context.Context,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "RequestPairingStage2")
	defer func() { logger.Tracef(ctx, "/RequestPairingStage2: %v", _err) }()
	msg := s.GetMessagePairingStage2()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToWiFiGroundStation) GetMessagePairingStage2() *duml.Message {
	return &duml.Message{
		Interface: duml.InterfaceIDAppToPairer,
		ID:        duml.MessageIDPairingStage2,
		Type:      duml.MessageTypePairingStage2,
		Payload:   []byte{0x31, 0x31, 0x00, 0x00, 0x00},
	}
}
