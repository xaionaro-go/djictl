package djible

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xaionaro-go/djictl/pkg/duml"
	"github.com/xaionaro-go/gatt"
)

type mockPeripheral struct {
	gatt.Peripheral
	writeFunc func(c *gatt.Characteristic, b []byte, noResp bool) error
}

func (m *mockPeripheral) WriteCharacteristic(ctx context.Context, c *gatt.Characteristic, b []byte, noResp bool) error {
	if m.writeFunc != nil {
		return m.writeFunc(c, b, noResp)
	}
	return nil
}

func TestDevice_SendMessage(t *testing.T) {
	mock := &mockPeripheral{}
	dev := NewDevice(mock, nil, duml.DeviceTypeOsmoAction4, "test-device")
	dev.CharacteristicSender = &gatt.Characteristic{}
	dev.CharacteristicReceiver = &gatt.Characteristic{}
	dev.CharacteristicPairingRequestor = &gatt.Characteristic{}

	ctx := context.Background()
	msg := &duml.Message{
		Interface: duml.InterfaceID{
			Sender:   duml.ComponentIDApp,
			Receiver: duml.ComponentIDCamera,
		},
		Type: duml.MessageType{
			CmdSet: duml.CommandSetCamera,
			CmdID:  duml.CommandID(0x01),
		},
	}

	writeCalled := false
	mock.writeFunc = func(c *gatt.Characteristic, b []byte, noResp bool) error {
		writeCalled = true
		return nil
	}

	err := dev.SendMessage(ctx, msg, true)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if !writeCalled {
		t.Error("WriteCharacteristic was not called")
	}
}

func TestDevice_Request_NoAckRequired(t *testing.T) {
	mock := &mockPeripheral{}
	dev := NewDevice(mock, nil, duml.DeviceTypeOsmoAction4, "test-device")

	ctx := context.Background()
	msg := &duml.Message{
		Type: duml.MessageType{
			Flags: 0, // No AckRequired
		},
	}

	_, err := dev.Request(ctx, msg, true)
	if err == nil {
		t.Fatal("Request should have failed for message without AckRequired flag")
	}
}

func TestDevice_Request_Success(t *testing.T) {
	mock := &mockPeripheral{}
	dev := NewDevice(mock, nil, duml.DeviceTypeOsmoAction4, "test-device")
	senderChar := &gatt.Characteristic{}
	dev.CharacteristicSender = senderChar
	dev.CharacteristicReceiver = &gatt.Characteristic{}
	dev.CharacteristicPairingRequestor = &gatt.Characteristic{}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	msg := &duml.Message{
		Interface: duml.InterfaceID{
			Sender:   duml.ComponentIDApp,
			Receiver: duml.ComponentIDCamera,
		},
		ID: 123,
		Type: duml.MessageType{
			Flags:  duml.MessageTypeFlagAckRequired,
			CmdSet: duml.CommandSetCamera,
			CmdID:  0x01,
		},
	}

	// Simulate receiving a response in a separate goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		resp := &duml.Message{
			Interface: duml.InterfaceID{
				Sender:   duml.ComponentIDCamera,
				Receiver: duml.ComponentIDApp,
			},
			ID: 123,
			Type: duml.MessageType{
				Flags:  duml.MessageTypeFlagResponse,
				CmdSet: duml.CommandSetCamera,
				CmdID:  0x01,
			},
		}
		dev.receiveNotification(context.Background(), senderChar, resp.Bytes(), nil)
	}()

	resp, err := dev.Request(ctx, msg, true)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.ID != 123 {
		t.Errorf("Expected ID 123, got %d", resp.ID)
	}
	if resp.Type.Flags&duml.MessageTypeFlagResponse == 0 {
		t.Error("Expected response flag to be set")
	}
}

func TestDevice_Request_Timeout(t *testing.T) {
	mock := &mockPeripheral{}
	dev := NewDevice(mock, nil, duml.DeviceTypeOsmoAction4, "test-device")
	dev.CharacteristicSender = &gatt.Characteristic{}
	dev.CharacteristicReceiver = &gatt.Characteristic{}
	dev.CharacteristicPairingRequestor = &gatt.Characteristic{}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	msg := &duml.Message{
		ID: 456,
		Type: duml.MessageType{
			Flags:  duml.MessageTypeFlagAckRequired,
			CmdSet: duml.CommandSetCamera,
			CmdID:  0x01,
		},
	}

	_, err := dev.Request(ctx, msg, true)
	if err == nil {
		t.Fatal("Request should have timed out")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got %v", err)
	}
}
