package pool

import (
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"sync"
)

const (
	SEND_BUFFER = 24
	RECV_BUFFER = 24
	TIMEOUT_SEC = 5
)

type Pool struct {
	lock       sync.RWMutex
	nodeID     uint32
	timeoutSec uint
	sendBuffer chan block.Block
	recvBuffer chan block.Block
	Manager    Manager
	Tunnels    []Worker
}

func NewDefaultPool(nodeID uint32) *Pool {
	pool := &Pool{
		nodeID:     nodeID,
		timeoutSec: TIMEOUT_SEC,
		sendBuffer: make(chan block.Block, SEND_BUFFER),
		recvBuffer: make(chan block.Block, RECV_BUFFER),
		Tunnels:    make([]Worker, 0),
	}
	return pool
}

func NewClientPool(dest string, nodeID uint32, cipher tunnel.Cipher, createWaitSec []uint) *Pool {
	pool := NewDefaultPool(nodeID)
	pool.Manager = NewClientManager(pool, dest, cipher, createWaitSec)
	go pool.Manager.Daemon()
	return pool
}

func NewServerPool(nodeID uint32, cipher tunnel.Cipher) *Pool {
	pool := NewDefaultPool(nodeID)
	pool.Manager = NewServerManager(pool, cipher)
	go pool.Manager.Daemon()
	return pool
}

func NewPreparingPool(cipher tunnel.Cipher, pushFunc func(clientID uint32, worker Worker)) *Pool {
	pool := NewDefaultPool(0)
	pool.Manager = NewPreparingManager(pool, cipher, pushFunc)
	go pool.Manager.Daemon()
	return pool
}

func (p *Pool) SendBlock(block block.Block) {
	p.sendBuffer <- block
}

func (p *Pool) RecvBlock() block.Block {
	return <-p.recvBuffer
}

func (p *Pool) RegisterWorker(worker Worker) {
	p.lock.Lock()
	p.Tunnels = append(p.Tunnels, worker)
	worker.StartRecv(p.recvBuffer)
	worker.StartSend(p.sendBuffer)
	p.lock.Unlock()
}

func (p *Pool) UnregisterWorker(worker Worker) {
	p.lock.Lock()
	for i, v := range p.Tunnels {
		if v == worker {
			p.Tunnels[i] = p.Tunnels[len(p.Tunnels)-1]
			p.Tunnels = p.Tunnels[:len(p.Tunnels)-1]
			break
		}
	}
	worker.Exit()
	p.lock.Unlock()
}

func (p *Pool) GetTunnelSize() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return len(p.Tunnels)
}
