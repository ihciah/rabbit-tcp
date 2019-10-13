package peer

import (
	"context"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"net"
	"sync"
)

type PeerGroup struct {
	lock        sync.Mutex
	cipher      tunnel.Cipher
	peerMapping map[uint32]*ServerPeer
	logger      *logger.Logger
}

func NewPeerGroup(cipher tunnel.Cipher) PeerGroup {
	if initRand() != nil {
		panic("Error when initialize random seed.")
	}
	return PeerGroup{
		cipher:      cipher,
		peerMapping: make(map[uint32]*ServerPeer),
		logger:      logger.NewLogger("[PeerGroup]"),
	}
}

// Add a tunnel to it's peer; will create peer if not exists
func (pg *PeerGroup) AddTunnel(tunnel *tunnel_pool.Tunnel) error {
	// add tunnel to peer(if absent, create peer to peer_group)
	pg.lock.Lock()
	var peer *ServerPeer
	var ok bool

	peerID := tunnel.GetPeerID()
	if peer, ok = pg.peerMapping[peerID]; !ok {
		peerContext, removePeerFunc := context.WithCancel(context.Background())
		serverPeer := NewServerPeerWithID(peerID, peerContext, removePeerFunc)
		peer = &serverPeer
		pg.peerMapping[peerID] = peer
		pg.logger.Infof("Server Peer %d added to PeerGroup.\n", peerID)

		go func() {
			<-peerContext.Done()
			pg.RemovePeer(peerID)
		}()
	}
	pg.lock.Unlock()
	peer.tunnelPool.AddTunnel(tunnel)
	return nil
}

// Like AddTunnel, add a raw connection
func (pg *PeerGroup) AddTunnelFromConn(conn net.Conn) error {
	tun, err := tunnel_pool.NewPassiveTunnel(conn, pg.cipher)
	if err != nil {
		conn.Close()
		return err
	}
	return pg.AddTunnel(&tun)
}

func (pg *PeerGroup) RemovePeer(peerID uint32) {
	pg.logger.Infof("Server Peer %d removed from peer group.\n", peerID)
	pg.lock.Lock()
	defer pg.lock.Unlock()
	delete(pg.peerMapping, peerID)
}
