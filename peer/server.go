package peer

import (
	"github.com/ihciah/rabbit-tcp/connection_pool"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
)

type ServerPeer struct {
	Peer
}

func NewServerPeerWithID(peerID uint32) ServerPeer {
	poolManager := tunnel_pool.NewServerManager()
	tunnelPool := tunnel_pool.NewTunnelPool(peerID, &poolManager)

	connectionPool := connection_pool.NewConnectionPool(&tunnelPool)

	return ServerPeer{
		Peer: Peer{
			peerID:         peerID,
			connectionPool: connectionPool,
			tunnelPool:     tunnelPool,
		},
	}
}
