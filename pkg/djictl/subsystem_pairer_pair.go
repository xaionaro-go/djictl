package djictl

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
)

const (
	defaultPINCode = "5160"
)

func (s *SubsystemPairer) Pair(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "Pair")
	defer func() { logger.Tracef(ctx, "/Pair: %v", _err) }()
	err := s.SendRequestStartPairing(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the request to start pairing: %w", err)
	}
	err = s.SendMessageSetPairingPIN(ctx, defaultPINCode)
	if err != nil {
		return fmt.Errorf("unable to send the message to set the PIN: %w", err)
	}

	logger.Debugf(ctx, "waiting for the pairing info")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, MessageTypePairingStatus):
		if len(msg.Payload) < 2 {
			logger.Errorf(ctx, "the payload size of the pairing status is not 4: %d", len(msg.Payload))
		} else {
			if msg.Payload[1] == 0x01 {
				logger.Debugf(ctx, "is already paired")
				return nil
			}
		}
	}

	logger.Debugf(ctx, "waiting for PIN approve")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, MessageTypePairingPINApproved):
		logger.Debugf(ctx, "PIN was approved: %#+v", msg)
	}

	err = s.SendMessagePairingStage1(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the command to pair (stage1): %w", err)
	}
	err = s.SendMessagePairingStage2(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the command to pair (stage2): %w", err)
	}
	return nil
}

func (s *SubsystemPairer) SendRequestStartPairing(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendRequestStartPairing")
	defer func() { logger.Tracef(ctx, "/SendRequestStartPairing: %v", _err) }()
	return s.Device().SendPairingRequest(ctx)
}

func (s *SubsystemPairer) SendMessageSetPairingPIN(
	ctx context.Context,
	pinCode string,
) (_err error) {
	logger.Tracef(ctx, "SendMessageSetPairingPIN")
	defer func() { logger.Tracef(ctx, "/SendMessageSetPairingPIN: %v", _err) }()
	msg := s.GetMessageSetPairingPIN(pinCode)
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemPairer) GetMessageSetPairingPIN(
	pinCode string,
) *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDSetPairingPIN,
		Type:      MessageTypeSetPairingPIN,
		Payload: s.GetMessagePayloadSetPairingPIN(
			pinCode,
		),
	}
}

func (s *SubsystemPairer) GetMessagePayloadSetPairingPIN(
	pinCode string,
) []byte {
	var buf bytes.Buffer
	must(buf.Write(packString("001749319286102")))
	must(buf.Write(packString(pinCode)))
	return buf.Bytes()
}

func (s *SubsystemPairer) ReceiveMessageSetPairingPINResult(
	ctx context.Context,
) (*Message, error) {
	return s.Device().ReceiveMessage(ctx, MessageTypeSetPairingPIN)
}

func (s *SubsystemPairer) SendMessagePairingStage1(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendMessagePairingStage1")
	defer func() { logger.Tracef(ctx, "/SendMessagePairingStage1: %v", _err) }()
	msg := s.GetMessagePairingStage1()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemPairer) GetMessagePairingStage1() *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDPairingStage1,
		Type:      MessageTypePairingStage1,
		Payload:   []byte{0x00},
	}
}

func (s *SubsystemPairer) SendMessagePairingStage2(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendMessagePairingStage2")
	defer func() { logger.Tracef(ctx, "/SendMessagePairingStage2: %v", _err) }()
	msg := s.GetMessagePairingStage2()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemPairer) GetMessagePairingStage2() *Message {
	return &Message{
		Subsystem: SubsystemIDOneMorePairer,
		ID:        MessageIDPairingStage2,
		Type:      MessageTypePairingStage2,
		Payload:   []byte{0x31, 0x31, 0x00, 0x00, 0x00},
	}
}
