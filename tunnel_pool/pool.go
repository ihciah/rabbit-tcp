package tunnel_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"sync"
)

type TunnelPool struct {
	mutex          sync.Mutex
	tunnelMapping  map[uint32]*Tunnel
	peerID         uint32
	manager        Manager
	sendQueue      chan block.Block
	sendRetryQueue chan block.Block
	recvQueue      chan block.Block
	ctx            context.Context
	cancel         context.CancelFunc // currently useless
	logger         *logger.Logger
}

func NewTunnelPool(peerID uint32, manager Manager, peerContext context.Context) *TunnelPool {
	ctx, cancel := context.WithCancel(peerContext)
	tp := &TunnelPool{
		tunnelMapping:  make(map[uint32]*Tunnel),
		peerID:         peerID,
		manager:        manager,
		sendQueue:      make(chan block.Block, SendQueueSize),
		sendRetryQueue: make(chan block.Block, SendQueueSize),
		recvQueue:      make(chan block.Block, RecvQueueSize),
		ctx:            ctx,
		cancel:         cancel,
		logger:         logger.NewLogger("[TunnelPool]"),
	}
	tp.logger.Infof("Tunnel Pool of peer %d created.\n", peerID)
	go manager.DecreaseNotify(tp)
	return tp
}

// Add a tunnel to tunnelPool and start bi-relay
func (tp *TunnelPool) AddTunnel(tunnel *Tunnel) {
	tp.logger.Infof("Tunnel %d added to Peer %d.\n", tunnel.tunnelID, tp.peerID)
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.tunnelMapping[tunnel.tunnelID] = tunnel
	tp.manager.Notify(tp)

	tunnel.ctx, tunnel.cancel = context.WithCancel(tp.ctx)
	go func() {
		<-tunnel.ctx.Done()
		tp.RemoveTunnel(tunnel)
	}()

	go tunnel.OutboundRelay(tp.sendQueue, tp.sendRetryQueue)
	go tunnel.InboundRelay(tp.recvQueue)
}

// Remove a tunnel from tunnelPool and stop bi-relay
func (tp *TunnelPool) RemoveTunnel(tunnel *Tunnel) {
	tp.logger.Infof("Tunnel %d to peer %d removed from pool.\n", tunnel.tunnelID, tunnel.peerID)
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	if tunnel, ok := tp.tunnelMapping[tunnel.tunnelID]; ok {
		delete(tp.tunnelMapping, tunnel.tunnelID)
		tp.manager.Notify(tp)
		go tp.manager.DecreaseNotify(tp)
	}
}

func (tp *TunnelPool) GetSendQueue() chan block.Block {
	return tp.sendQueue
}

func (tp *TunnelPool) GetRecvQueue() chan block.Block {
	return tp.recvQueue
}
