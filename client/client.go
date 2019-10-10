package client

import (
	"context"
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	peer   peer.ClientPeer
	logger *log.Logger
}

func NewClient(tunnelNum int, endpoint string, cipher tunnel.Cipher) Client {
	return Client{
		peer:   peer.NewClientPeer(tunnelNum, endpoint, cipher),
		logger: log.New(os.Stdout, "[Client]", log.LstdFlags),
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
			c.logger.Printf("Error when accept connection: %v.\n", err)
		}
		connProxy := c.Dial(dest)
		biRelay(conn, connProxy)
	}
}

func biRelay(left, right net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	go relay(left, right, ctx, cancel)
	go relay(right, left, ctx, cancel)
	<-ctx.Done()
	left.Close()
	right.Close()
}

func relay(dst, src net.Conn, ctx context.Context, cancel context.CancelFunc) {
	_, err := io.Copy(dst, src)
	if err != nil {
		cancel()
		dst.SetDeadline(time.Now())
		src.SetDeadline(time.Now())
	}
}
