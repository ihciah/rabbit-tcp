package client

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
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

func (c *Client) Dial(address string) connection.HalfOpenConn {
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
		go func() {
			c.logger.Infoln("Accepted a connection.")
			connProxy := c.Dial(dest)
			biRelay(conn.(*net.TCPConn), connProxy, c.logger)
		}()
	}
}

func biRelay(left, right connection.HalfOpenConn, logger *logger.Logger) {
	var wg sync.WaitGroup
	wg.Add(1)
	go relay(left, right, &wg, logger, "local <- tunnel")
	wg.Add(1)
	go relay(right, left, &wg, logger, "local -> tunnel")
	wg.Wait()
	// logger.Errorf("===========> Close client biRelay")
	_ = left.Close()
	_ = right.Close()
}

func relay(dst, src connection.HalfOpenConn, wg *sync.WaitGroup, logger *logger.Logger, label string) {
	defer wg.Done()
	_, err := io.Copy(dst, src)
	if err != nil {
		_ = dst.SetDeadline(time.Now())
		_ = src.SetDeadline(time.Now())
		_ = dst.Close()
		_ = src.Close()
		if err != io.EOF {
			logger.Errorf("Error when relay client: %v.\n", err)
		}
	} else {
		// logger.Debugf("!!!!!!!!!!!!!!!! %s : dst close write", label)
		dst.CloseWrite()
		// logger.Debugf("!!!!!!!!!!!!!!!! %s : src close read", label)
		src.CloseRead()
	}
}
