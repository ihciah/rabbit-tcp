package client

import "github.com/ihciah/rabbit-tcp/peer"

type Client struct {
	//config
	peer peer.Peer
}

func (c *Client) Dial(address string) {
	c.peer.AddConnection()
}
