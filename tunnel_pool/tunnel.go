package tunnel_pool

import (
	"bytes"
	"context"
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

// Create a new tunnel from a net.Conn and cipher with random connectionID
func NewTunnel(conn net.Conn, ciph tunnel.Cipher) *Tunnel {
	ctx, cancel := context.WithCancel(context.Background())
	return &Tunnel{
		Conn:     tunnel.NewEncryptedConn(conn, ciph),
		ctx:      ctx,
		cancel:   cancel,
		tunnelID: rand.Uint32(),
	}
}

// Read block from send channel, pack it and send
func (tunnel *Tunnel) OutboundRelay(input <-chan block.Block) {
	for {
		select {
		case <-tunnel.ctx.Done():
			return
		case blk := <-input:
			reader := bytes.NewReader(blk.Pack())
			io.Copy(tunnel.Conn, reader)
			// TODO: error handle
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
