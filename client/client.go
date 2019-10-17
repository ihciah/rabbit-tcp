package client

import (
	"context"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"io"
	"net"
	"time"
)

type Client struct {
	peer   peer.ClientPeer
	logger *logger.Logger
}

func NewClient(tunnelNum int, endpoint string, cipher tunnel.Cipher) Client {
	return Client{
		peer:   peer.NewClientPeer(tunnelNum, endpoint, cipher),
		logger: logger.NewLogger("[Client]"),
	}
}

func (c *Client) Dial(address string) net.Conn {
	return c.peer.Dial(address)
}

func (c *Client) ServeForward(listen, dest string) error {
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			c.logger.Errorf("Error when accept connection: %v.\n", err)
			continue
		}
		c.logger.Infoln("Accepted a connection.")
		connProxy := c.Dial(dest)
		go biRelay(conn, connProxy, c.logger)
	}
}

func biRelay(left, right net.Conn, logger *logger.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	go relay(left, right, cancel, logger)
	go relay(right, left, cancel, logger)
	<-ctx.Done()
	_ = left.Close()
	_ = right.Close()
}

func relay(dst, src net.Conn, cancel context.CancelFunc, logger *logger.Logger) {
	_, err := io.Copy(dst, src)
	if err != nil {
		_ = dst.SetDeadline(time.Now())
		_ = src.SetDeadline(time.Now())
		cancel()
		if err != io.EOF {
			logger.Errorf("Error when relay client: %v.\n", err)
		}
	}
}
