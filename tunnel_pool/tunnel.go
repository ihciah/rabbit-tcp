package tunnel_pool

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"io"
	"math/rand"
	"net"
)

type Tunnel struct {
	net.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	tunnelID uint32
}

// Create a new tunnel from a net.Conn and cipher with random tunnelID
func NewTunnel(conn net.Conn, ciph tunnel.Cipher) (*Tunnel, error) {
	return NewTunnelWithID(conn, ciph, rand.Uint32())
}

// Create a new tunnel from a net.Conn and cipher with given tunnelID
func NewTunnelWithID(conn net.Conn, ciph tunnel.Cipher, tunnelID uint32) (*Tunnel, error) {
	ctx, cancel := context.WithCancel(context.Background())
	tun := &Tunnel{
		Conn:     tunnel.NewEncryptedConn(conn, ciph),
		ctx:      ctx,
		cancel:   cancel,
		tunnelID: tunnelID,
	}
	tunnelIDBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(tunnelIDBuffer, tun.tunnelID)
	_, err := io.Copy(tun.Conn, bytes.NewReader(tunnelIDBuffer))
	return tun, err
}

// Read block from send channel, pack it and send
func (tunnel *Tunnel) OutboundRelay(input <-chan block.Block) {
	for {
		select {
		case <-tunnel.ctx.Done():
			return
		case blk := <-input:
			reader := bytes.NewReader(blk.Pack())
			_, err := io.Copy(tunnel.Conn, reader)
			if err != nil {
				// TODO: error handle
			}
		}
	}
}

// Read bytes from connection, parse it to block then put in recv channel
func (tunnel *Tunnel) InboundRelay(output chan<- block.Block) {
	for {
		select {
		case <-tunnel.ctx.Done():
			return
		default:
			blk := block.NewBlockFromReader(tunnel.Conn)
			output <- *blk
		}
	}
}

func (tunnel *Tunnel) StopRelay() {
	tunnel.cancel()
}
