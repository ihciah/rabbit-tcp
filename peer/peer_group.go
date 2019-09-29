package peer

import (
	"encoding/binary"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"io"
	"sync"
)

type PeerGroup struct {
	lock        sync.Mutex
	peerMapping map[uint32]*ServerPeer
}

func (pg *PeerGroup) AddTunnel(tunnel tunnel_pool.Tunnel) error {
	// 1. read peer id
	// 2. add tunnel to peer(if absent, create peer to peer_group)
	peerIDBuffer := make([]byte, 4)
	_, err := io.ReadFull(tunnel, peerIDBuffer)
	peerID := binary.LittleEndian.Uint32(peerIDBuffer)

	if err != nil {
		return err
	}
	pg.lock.Lock()
	var peer *ServerPeer
	var ok bool

	if peer, ok = pg.peerMapping[peerID]; !ok {
		peer = NewServerPeer()
		pg.peerMapping[peerID] = peer
	}
	pg.lock.Unlock()
	peer.AddTunnel(tunnel)
	return nil
}

// TODO: if all tunnels down, after WAIT time, remove peer
func (pg *PeerGroup) RemovePeer(peerID uint32) {

}
