// Package main provides the entry point for the djictl command-line tool.
// This file contains unit tests for the djictl package using a simulated BLE device.
package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/djible"
	"github.com/xaionaro-go/djictl/pkg/duml"
	"github.com/xaionaro-go/gatt"
)

func TestConnectWiFiAndStartStreaming(t *testing.T) {
	ctx := getContext(logger.LevelDebug, false, "")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 1. Setup simulated device
	svc := gatt.NewService(gatt.MustParseUUID("0000180a-0000-1000-8000-00805f9b34fb"))
	simDevice := gatt.NewSimDeviceClient(svc, "DJI Osmo Pocket 3")
	simDevice.SetManufacturerData([]byte{0xAA, 0x08, 0x20, 0x00})

	// Receiver characteristic (0x002D)
	receiverChar := svc.AddCharacteristic(gatt.MustParseUUID("0000002d-0000-1000-8000-00805f9b34fb"))
	receiverChar.SetVHandle(0x002D)
	receiverChar.HandleNotifyFunc(func(ctx context.Context, r gatt.Request, n gatt.Notifier) {
		t.Log("Central subscribed to receiver")
		msg := &duml.Message{
			Type:    duml.MessageTypeMaybeStatus,
			Payload: make([]byte, 100), // Big enough packet
		}
		b := msg.Bytes()
		t.Logf("RECV_HEX: %X", b)
		simDevice.SendNotification(0x002D, b)
	})

	// Sender characteristic (0x0030)
	senderChar := svc.AddCharacteristic(gatt.MustParseUUID("00000030-0000-1000-8000-00805f9b34fb"))
	senderChar.SetVHandle(0x0030)
	senderChar.HandleWriteFunc(func(ctx context.Context, r gatt.Request, b []byte) (status byte) {
		t.Logf("SENT_HEX: %X", b)
		msg, err := duml.ParseMessage(b)
		if err != nil {
			t.Errorf("failed to parse message: %v", err)
			return 0
		}
		t.Logf("Device received message: %v", msg.Type)

		switch msg.Type {
		case duml.MessageTypePrepareToLiveStream:
			go func() {
				resp := &duml.Message{
					Interface: msg.Interface,
					ID:        msg.ID,
					Type:      duml.MessageTypePrepareToLiveStreamResult,
					Payload:   []byte{0x00}, // Success
				}
				b := resp.Bytes()
				t.Logf("RECV_HEX: %X", b)
				simDevice.SendNotification(0x002D, b)
			}()
		case duml.MessageTypeConnectToWiFi:
			go func() {
				resp := &duml.Message{
					Interface: msg.Interface,
					ID:        msg.ID,
					Type:      duml.MessageTypeConnectToWiFiResult,
					Payload:   []byte{0x00, 0x00}, // Success
				}
				b := resp.Bytes()
				t.Logf("RECV_HEX: %X", b)
				simDevice.SendNotification(0x002D, b)
			}()
		case duml.MessageTypeStartStopStreaming: // Also MessageTypeConfigure
			go func() {
				resp := &duml.Message{
					Interface: msg.Interface,
					ID:        msg.ID,
					Type:      msg.Type,     // Use the same type for response (or the result type if known)
					Payload:   []byte{0x00}, // Success
				}
				// For StartStopStreaming, the result type is actually different
				if msg.Interface == duml.InterfaceIDAppToVideoTransmission {
					resp.Type = duml.MessageTypeStartStopStreamingResult
				}
				b := resp.Bytes()
				t.Logf("RECV_HEX: %X", b)
				simDevice.SendNotification(0x002D, b)

				// If this is the actual start streaming message (not the prepare stage),
				// then send the streaming status and cancel the context.
				if msg.Interface == duml.InterfaceIDAppToVideoTransmission && len(msg.Payload) > 4 && msg.Payload[0] == 0x01 {
					// Send streaming status to satisfy the loop in LiveStream
					status := &duml.Message{
						Type:    duml.MessageTypeBatteryStatus,
						Payload: make([]byte, 21),
					}
					status.Payload[20] = 100 // 100% battery
					b := status.Bytes()
					t.Logf("RECV_HEX: %X", b)
					simDevice.SendNotification(0x002D, b)

					// Cancel the context to stop the infinite loop in LiveStream
					cancel()
				}
			}()
		case duml.MessageTypeSetPairingPIN:
			go func() {
				resp := &duml.Message{
					Interface: msg.Interface,
					ID:        msg.ID,
					Type:      duml.MessageTypePairingStatus,
					Payload:   []byte{0x00, 0x01}, // Already paired
				}
				b := resp.Bytes()
				t.Logf("RECV_HEX: %X", b)
				simDevice.SendNotification(0x002D, b)
			}()
		}
		return 0
	})

	// Pairing Requestor characteristic (0x002E)
	pairingChar := svc.AddCharacteristic(gatt.MustParseUUID("0000002e-0000-1000-8000-00805f9b34fb"))
	pairingChar.SetVHandle(0x002E)
	pairingChar.HandleWriteFunc(func(ctx context.Context, r gatt.Request, b []byte) (status byte) {
		t.Logf("SENT_HEX: %X", b)
		return 0
	})

	// 2. Start scanning
	devCh, errCh, err := djible.ScanWithDevice(ctx, simDevice)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	select {
	case dev := <-devCh:
		t.Logf("Found device: %s", dev)

		err := connectWiFiAndStartStreaming(ctx, dev, "test-ssid", "test-psk", "rtmp://test/live", duml.Resolution1080p, 6000, duml.FPS30)
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("runProcess failed: %v", err)
		}

	case err := <-errCh:
		t.Fatalf("Error during scan: %v", err)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for device")
	}
}
