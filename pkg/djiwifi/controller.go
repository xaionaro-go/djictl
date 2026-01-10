package djiwifi

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

type Controller struct {
	conn *net.UDPConn
	addr *net.UDPAddr

	nextSeq atomic.Uint32

	ctx    context.Context
	cancel context.CancelFunc
}

func NewController(ctx context.Context, deviceAddr string) (*Controller, error) {
	addr, err := net.ResolveUDPAddr(ProtocolUDP, deviceAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve UDP address %s: %w", deviceAddr, err)
	}

	conn, err := net.DialUDP(ProtocolUDP, nil, addr)
	if err != nil {
		return nil, fmt.Errorf("unable to dial UDP to %s: %w", deviceAddr, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	c := &Controller{
		conn:   conn,
		addr:   addr,
		ctx:    ctx,
		cancel: cancel,
	}

	return c, nil
}

func (c *Controller) Close() error {
	c.cancel()
	return c.conn.Close()
}

func (c *Controller) SendPacket(ctx context.Context, p *Packet) error {
	logger.Tracef(ctx, "SendPacket: Type=%s Len=%d", p.Type, len(p.Payload))
	_, err := c.conn.Write(p.Bytes())
	return err
}

func (c *Controller) SendDUML(ctx context.Context, msg *duml.Message, metadata Metadata) error {
	p := NewDUMLPacket(msg, metadata)
	return c.SendPacket(ctx, p)
}

func (c *Controller) SendAppStatus(ctx context.Context, metadata Metadata, payload []byte) error {
	p := &Packet{
		Type:     MessageTypeControl,
		Metadata: metadata,
		Payload:  payload,
	}
	return c.SendPacket(ctx, p)
}

func (c *Controller) ReceivePacket(ctx context.Context) (*Packet, error) {
	buf := make([]byte, ReadBufferSize)
	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetReadDeadline(deadline)
		defer c.conn.SetReadDeadline(time.Time{})
	}
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return ParsePacket(buf[:n])
}

func (c *Controller) SendHandshake(ctx context.Context) error {
	if err := c.SendAppStatus(ctx, MetadataInitial, PayloadInitial); err != nil {
		return fmt.Errorf("failed to send initial status: %w", err)
	}

	msgApp := &duml.Message{
		Interface: duml.InterfaceIDAppToApp,
		ID:        duml.MessageIDAppIdentifier,
		Payload:   PayloadAppIdentifier,
	}
	if err := c.SendDUML(ctx, msgApp, MetadataApp); err != nil {
		return fmt.Errorf("failed to send sAPP command: %w", err)
	}

	return nil
}

// SendVideoHandshake sends the RMVT magic handshake to trigger video streaming.
func (c *Controller) SendVideoHandshake(ctx context.Context) error {
	p := &Packet{
		WhType:  WhTypeHandshake,
		Payload: PayloadHandshakeRMVT,
	}
	return c.SendPacket(ctx, p)
}

// SendSimulatorData sends Remote Controller stick and button data in simulator mode.
func (c *Controller) SendSimulatorData(ctx context.Context, data duml.RemoteControllerSimulatorData) error {
	msg := duml.NewRemoteControllerSimulatorMessage(data)
	return c.SendDUML(ctx, msg, MetadataApp)
}

func (c *Controller) SendFCCEnable(ctx context.Context) error {
	msg := duml.NewFCCEnableMessage()
	msg.Interface = duml.InterfaceIDAppToCamera
	return c.SendDUML(ctx, msg, MetadataApp)
}

func (c *Controller) SendConfigureBroadcast(ctx context.Context, url string, enable bool) error {
	msg := duml.NewBroadcastMessage(enable, url)
	return c.SendDUML(ctx, msg, MetadataApp)
}

// SendStopStreaming sends the command to stop video streaming.
func (c *Controller) SendStopStreaming(ctx context.Context) error {
	msg := &duml.Message{
		Interface: duml.InterfaceIDAppToVideoTransmission,
		ID:        duml.MessageIDStopStreaming,
		Type:      duml.MessageTypeStartStopStreaming,
		Payload:   []byte{0x00}, // 0 = Stop
	}
	return c.SendDUML(ctx, msg, MetadataApp)
}
