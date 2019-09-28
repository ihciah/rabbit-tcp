package peer

import (
	"github.com/ihciah/rabbit-tcp/connection_pool"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
)

type Peer struct {
	peerID         uint32
	connectionPool connection_pool.ConnectionPool
	tunnelPool     tunnel_pool.TunnelPool
}

func (p *Peer) AddConnection() {

}

func (p *Peer) AddTunnel() {

}
