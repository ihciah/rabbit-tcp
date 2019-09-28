package peer

import "net"

type PeerGroup struct {
	peerMapping map[uint32]*ServerPeer
}

func (pg *PeerGroup) AddTunnel(conn net.Conn) {
	// 1. read peer id
	// 2. add tunnel to peer(if absent, create peer to peer_group)
}

// TODO: if all tunnels down, after WAIT time, remove peer
func (pg *PeerGroup) RemovePeer(peerID uint32) {

}
