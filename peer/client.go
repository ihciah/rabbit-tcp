package peer

import (
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/connection_pool"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"math/rand"
)

type ClientPeer struct {
	Peer
}

func NewClientPeer(tunnelNum int, endpoint string, cipher tunnel.Cipher) ClientPeer {
	peerID := rand.Uint32()
	return NewClientPeerWithID(peerID, tunnelNum, endpoint, cipher)
}

func NewClientPeerWithID(peerID uint32, tunnelNum int, endpoint string, cipher tunnel.Cipher) ClientPeer {
	poolManager := tunnel_pool.NewClientManager(tunnelNum, endpoint, peerID, cipher)
	tunnelPool := tunnel_pool.NewTunnelPool(peerID, &poolManager)

	connectionPool := connection_pool.NewConnectionPool(&tunnelPool)

	return ClientPeer{
		Peer: Peer{
			peerID:         peerID,
			connectionPool: connectionPool,
			tunnelPool:     tunnelPool,
		},
	}
}

func (cp *ClientPeer) Dial(address string) connection.Connection {
	conn := cp.connectionPool.NewPooledInboundConnection()
	conn.SendConnect(address)
	return conn
}
