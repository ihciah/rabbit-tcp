package tunnel_pool

import (
	"github.com/ihciah/rabbit-tcp/block"
	"sync"
)

type TunnelPool struct {
	mutex         sync.Mutex
	tunnelMapping map[uint32]*Tunnel
	manager       Manager
	sendQueue     chan block.Block
	recvQueue     chan block.Block
}

func NewTunnelPool(manager Manager) {
	// TODO: remember call notify at create
}

// Add a tunnel to tunnelPool and start bi-relay
func (tp *TunnelPool) AddTunnel(tunnel *Tunnel) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.tunnelMapping[tunnel.tunnelID] = tunnel
	tp.manager.Notify(tp)
	go tunnel.OutboundRelay(tp.sendQueue)
	go tunnel.InboundRelay(tp.recvQueue)
}

// Remove a tunnel from tunnelPool and stop bi-relay
func (tp *TunnelPool) RemoveTunnel(tunnel *Tunnel) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	if tunnel, ok := tp.tunnelMapping[tunnel.tunnelID]; ok {
		tunnel.StopRelay()
		delete(tp.tunnelMapping, tunnel.tunnelID)
		tp.manager.Notify(tp)
		tp.manager.DecreaseNotify(tp)
	}
}

func (tp *TunnelPool) GetTunnelCount() int {
	return len(tp.tunnelMapping)
}

func (tp *TunnelPool) GetSendQueue() chan block.Block {
	return tp.sendQueue
}

func (tp *TunnelPool) GetRecvQueue() chan block.Block {
	return tp.recvQueue
}
