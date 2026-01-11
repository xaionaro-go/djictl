package djible

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToVideoTransmission) StopLiveStream(
	ctx context.Context,
) error {
	_, err := s.RequestStopLiveStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the duml.Message: %w", err)
	}
	return nil
}

func (s *InterfaceAppToVideoTransmission) RequestStopLiveStream(
	ctx context.Context,
) (*duml.Message, error) {
	msg := s.GetMessageStopLiveStream()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToVideoTransmission) GetMessageStopLiveStream() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDStopStreaming,
		Type:      duml.MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadStopLiveStream(),
	}
}

func (s *InterfaceAppToVideoTransmission) GetMessagePayloadStopLiveStream() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01, 0x1A, 0x00, 0x01, 0x02}))
	return buf.Bytes()
}
