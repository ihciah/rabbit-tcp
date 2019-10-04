package tunnel_pool

import (
	"github.com/ihciah/rabbit-tcp/block"
	"log"
	"os"
	"sync"
)

const (
	SendQueueSize = 24
	RecvQueueSize = 24
)

type TunnelPool struct {
	mutex         sync.Mutex
	tunnelMapping map[uint32]*Tunnel
	peerID        uint32
	manager       Manager
	sendQueue     chan block.Block
	recvQueue     chan block.Block
	logger        *log.Logger
}

func NewTunnelPool(peerID uint32, manager Manager) TunnelPool {
	// TODO: remember call notify at create
	tp := TunnelPool{
		tunnelMapping: make(map[uint32]*Tunnel),
		peerID:        peerID,
		manager:       manager,
		sendQueue:     make(chan block.Block, SendQueueSize),
		recvQueue:     make(chan block.Block, RecvQueueSize),
		logger:        log.New(os.Stdout, "[TunnelPool]", log.LstdFlags),
	}
	go manager.DecreaseNotify(&tp)
	return tp
}

// Add a tunnel to tunnelPool and start bi-relay
func (tp *TunnelPool) AddTunnel(tunnel *Tunnel) {
	tp.logger.Println("AddTunnel called with tunnel", tunnel.tunnelID)

	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.tunnelMapping[tunnel.tunnelID] = tunnel
	tp.manager.Notify(tp)
	go tunnel.OutboundRelay(tp.sendQueue)
	go tunnel.InboundRelay(tp.recvQueue)
}

// Remove a tunnel from tunnelPool and stop bi-relay
func (tp *TunnelPool) RemoveTunnel(tunnel *Tunnel) {
	tp.logger.Println("RemoveTunnel called with tunnel", tunnel.tunnelID)

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
