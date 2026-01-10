package djible

import (
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
	"github.com/xaionaro-go/gatt"
	"github.com/xaionaro-go/xsync"
)

const (
	characteristicIDReceiver         = uint16(0x002D)
	characteristicIDPairingRequestor = uint16(0x002E)
	characteristicIDSender           = uint16(0x0030)
)

type Device struct {
	Periph gatt.Peripheral
	ID     DeviceID
	Type   duml.DeviceType
	Name   string

	ConnectedChan                  chan struct{}
	CharacteristicSender           *gatt.Characteristic
	CharacteristicPairingRequestor *gatt.Characteristic
	CharacteristicReceiver         *gatt.Characteristic

	ReceiveLocker                          xsync.Mutex
	ReceivedPairingRequestConfirmationChan chan struct{}
	ReceivedMessageChan                    map[duml.MessageType]chan *duml.Message
}

func NewDevice(
	periph gatt.Peripheral,
	id net.HardwareAddr,
	typ duml.DeviceType,
	name string,
) *Device {
	return &Device{
		Periph: periph,
		ID:     id,
		Type:   typ,
		Name:   name,

		ConnectedChan:                          make(chan struct{}),
		ReceivedPairingRequestConfirmationChan: make(chan struct{}),
		ReceivedMessageChan:                    make(map[duml.MessageType]chan *duml.Message),
	}
}

func (d *Device) String() string {
	return fmt.Sprintf("%s (%s)", d.ID, d.Name)
}

func (d *Device) getCharacteristics(
	ctx context.Context,
) (receiver *gatt.Characteristic, sender *gatt.Characteristic, pairingRequestor *gatt.Characteristic, _err error) {
	logger.Debugf(ctx, "discovering services...")

	services, err := d.Periph.DiscoverServices(ctx, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to discover services: %w", err)
	}
	logger.Debugf(ctx, "received %d services", len(services))

	sort.Slice(services, func(i, j int) bool {
		return services[i].Handle() < services[j].Handle()
	})

	for _, service := range services {
		logger.Debugf(ctx, "probing service %s (0x%02X:0x%X) for the magic characteristic...", service.Name(), service.Handle(), service.UUID().Bytes())
		characteristics, err := d.Periph.DiscoverCharacteristics(
			ctx,
			nil,
			service,
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to discover characteristics of service %s:%s: %w", service.UUID(), service.Name(), err)
		}
		logger.Debugf(ctx, "received %d characteristics", len(characteristics))
		for _, characteristic := range characteristics {
			logger.Tracef(ctx, "found characteristic %04X:%04X", characteristic.Handle(), characteristic.VHandle())
			switch characteristic.VHandle() {
			case characteristicIDReceiver:
				receiver = characteristic
			case characteristicIDSender:
				sender = characteristic
			case characteristicIDPairingRequestor:
				pairingRequestor = characteristic
			default:
				continue
			}
		}
	}

	switch {
	case receiver == nil:
		return nil, nil, nil, fmt.Errorf("unable to find characteristic %04X", characteristicIDReceiver)
	case sender == nil:
		return nil, nil, nil, fmt.Errorf("unable to find characteristic %04X", characteristicIDSender)
	case pairingRequestor == nil:
		pairingRequestor = gatt.NewCharacteristic(gatt.UUID{}, nil, 0, 0, characteristicIDPairingRequestor)
	}
	return
}

func (d *Device) Init(ctx context.Context) (_err error) {
	logger.Tracef(ctx, "Init(ctx)")
	defer func() { logger.Tracef(ctx, "/Init(ctx): %v %v", _err) }()

	logger.Debugf(ctx, "connecting to %s:%s", d.Periph.ID(), d.Periph.Name())
	d.Periph.Device().Connect(ctx, d.Periph)
	<-d.ConnectedChan
	d.Periph.Subscribe(0x2D, func(c *gatt.Characteristic, b []byte, err error) {
		d.receiveNotification(ctx, c, b, err)
	})
	logger.Debugf(ctx, "connected")

	d.Periph.DiscoverIncludedServices(ctx, nil, nil) // imitating DJI MIMO

	characteristicReceiver, characteristicSender, characteristicPairingRequestor, err := d.getCharacteristics(ctx)
	if err != nil {
		return fmt.Errorf("unable to get the magic characteristic: %w", err)
	}
	d.CharacteristicReceiver = characteristicReceiver
	d.CharacteristicPairingRequestor = characteristicPairingRequestor
	d.CharacteristicSender = characteristicSender
	logger.Debugf(ctx, "set the magic characteristics: %04X %04X", characteristicReceiver.VHandle(), characteristicSender.VHandle())

	d.Periph.DiscoverDescriptors(ctx, nil, &gatt.Characteristic{}) // imitating DJI MIMO

	logger.Debugf(ctx, "setting MTU to 517")
	if err := d.Periph.SetMTU(ctx, 517); err != nil {
		return fmt.Errorf("unable to set MTU: %w", err)
	}

	logger.Debugf(ctx, "waiting for the a big enough packet (to make sure MTU is already set to the correct value)")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-d.getReceiveMessageChan(ctx, duml.MessageTypeMaybeStatus):
		logger.Debugf(ctx, "received a status: %#+v", msg)
	}

	return nil
}

func (d *Device) receiveNotification(
	ctx context.Context,
	c *gatt.Characteristic,
	b []byte,
	err error,
) {
	logger.Tracef(ctx, "receiveNotification(ctx, %d:'%s', %X, %v)", c.VHandle(), c.Name(), b, err)
	defer func() {
		logger.Tracef(ctx, "/receiveNotification(ctx, %d:'%s', %X, %v)", c.VHandle(), c.Name(), b, err)
	}()

	if err != nil {
		logger.Errorf(ctx, "received a notification about an error: %v")
		return
	}

	msg, err := duml.ParseMessage(b)
	if err != nil {
		logger.Errorf(ctx, "unable to parse the duml.Message: %v", err)
		return
	}
	logger.Debugf(ctx, "received duml.Message: %#+v", msg)
	logger.Tracef(ctx, "payload: %X", msg.Payload)
	select {
	case d.getReceiveMessageChan(ctx, msg.Type) <- msg:
	default:
		logger.Debugf(ctx, "nobody waits for this message (%v), skipping", msg.Type)
	}
}

func (d *Device) getReceiveMessageChan(
	ctx context.Context,
	msgType duml.MessageType,
) chan *duml.Message {
	return xsync.DoR1(ctx, &d.ReceiveLocker, func() chan *duml.Message {
		if d.ReceivedMessageChan[msgType] == nil {
			d.ReceivedMessageChan[msgType] = make(chan *duml.Message, 1)
		}
		return d.ReceivedMessageChan[msgType]
	})
}

func (d *Device) IsInitialized() bool {
	return d.CharacteristicReceiver != nil && d.CharacteristicSender != nil && d.CharacteristicPairingRequestor != nil
}

func (d *Device) SendPairingRequest(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "SendPairingRequest")
	defer func() { logger.Tracef(ctx, "/SendPairingRequest: %v", _err) }()
	if !d.IsInitialized() {
		return fmt.Errorf("call Init first")
	}
	return d.Periph.WriteCharacteristic(ctx, d.CharacteristicPairingRequestor, []byte{0x01, 0x00}, false)
}

func (d *Device) SendMessage(
	ctx context.Context,
	msg *duml.Message,
	noResponse bool,
) (_err error) {
	logger.Tracef(ctx, "SendMessage")
	defer func() { logger.Tracef(ctx, "/SendMessage: %v", _err) }()
	if !d.IsInitialized() {
		return fmt.Errorf("call Init first")
	}
	return d.Periph.WriteCharacteristic(ctx, d.CharacteristicSender, msg.Bytes(), noResponse)
}

func (d *Device) ReceiveMessage(
	ctx context.Context,
	msgType duml.MessageType,
) (_ret *duml.Message, _err error) {
	logger.Tracef(ctx, "ReceiveMessage")
	defer func() { logger.Tracef(ctx, "/ReceiveMessage: %v", _err) }()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case v := <-d.getReceiveMessageChan(ctx, msgType):
		return v, nil
	}
}
