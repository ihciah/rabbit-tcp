package client

import (
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"net"
)

type Client struct {
	//config
	peer peer.ClientPeer
}

func NewClient(tunnelNum int, endpoint string, cipher tunnel.Cipher) Client {
	return Client{peer: peer.NewClientPeer(tunnelNum, endpoint, cipher)}
}

func (c *Client) Dial(address string) net.Conn {
	return c.peer.Dial(address)
}
