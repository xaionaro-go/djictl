package djible

import "github.com/xaionaro-go/djictl/pkg/duml"

import (
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/gatt"
	"github.com/xaionaro-go/xsync"
)

const (
	logUnknownDevices = false
)

func Scan(
	ctx context.Context,
) (<-chan *Device, <-chan error, error) {
	d, err := gatt.NewDevice(ctx,
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, true),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open device, err: %w", err)
	}
	return ScanWithDevice(ctx, d)
}

func ScanWithDevice(
	ctx context.Context,
	d gatt.Device,
) (<-chan *Device, <-chan error, error) {
	ctx, cancelFn := context.WithCancel(ctx)

	retCh := make(chan *Device, 100)
	errCh := make(chan error, 2)

	devices := map[string]*Device{}
	devicesLocker := xsync.Mutex{}

	d.Handle(
		ctx,
		gatt.PeripheralDiscovered(func(
			ctx context.Context,
			periph gatt.Peripheral,
			adv *gatt.Advertisement,
			rssi int,
		) {
			if logUnknownDevices {
				logger.Tracef(ctx, "gatt.PeripheralDiscovered(ctx, %s:%s)", periph.ID(), periph.Name())
				defer func() {
					logger.Tracef(ctx, "/gatt.PeripheralDiscovered(ctx, %s:%s)", periph.ID(), periph.Name())
				}()
			}
			deviceType := duml.IdentifyDeviceType(adv.ManufacturerData)
			if deviceType == duml.DeviceTypeUndefined {
				if logUnknownDevices {
					logger.Debugf(ctx, "ignoring device %s: considered a non DJI Osmo device (%X)", periph.ID(), adv.ManufacturerData)
				}
				return
			}
			deviceID, err := ParseDeviceID(periph.ID())
			if err != nil {
				cancelFn()
				errCh <- fmt.Errorf("unable to parse device ID '%s': %w", periph.ID(), err)
				return
			}
			dev := NewDevice(periph, deviceID, deviceType, adv.LocalName)
			devicesLocker.Do(ctx, func() {
				devices[periph.ID()] = dev
			})
			retCh <- dev
		}),
		gatt.PeripheralConnected(func(ctx context.Context, periph gatt.Peripheral, err error) {
			logger.Tracef(ctx, "gatt.PeripheralConnected(ctx, %s:%s, %v)", periph.ID(), periph.Name(), err)
			defer func() {
				logger.Tracef(ctx, "/gatt.PeripheralConnected(ctx, %s:%s, %v)", periph.ID(), periph.Name(), err)
			}()
			dev := xsync.DoR1(ctx, &devicesLocker, func() *Device {
				return devices[periph.ID()]
			})
			if dev == nil {
				logger.Errorf(ctx, "connected to unexpected device %s:%s", periph.ID(), periph.Name())
				return
			}
			dev.Periph = periph
			close(dev.ConnectedChan)
		}),
	)

	err := d.Start(ctx, func(ctx context.Context, d gatt.Device, s gatt.State) {
		logger.Debugf(ctx, "state changed on device %v to %v", d.ID(), s)
		switch s {
		case gatt.StatePoweredOn:
			err := d.Scan(ctx, nil, false)
			if err != nil {
				cancelFn()
				errCh <- fmt.Errorf("unable to start scanning: %w", err)
			}
			return
		default:
			cancelFn()
			errCh <- fmt.Errorf("received unexpected state: %s", s)
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize the bluetooth interface: %w", err)
	}

	return retCh, errCh, nil
}
