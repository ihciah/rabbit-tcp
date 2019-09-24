package connection

import (
	"github.com/ihciah/rabbit-tcp/pool"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"net"
	"sync"
)

type Collector struct {
	lock          sync.RWMutex
	mapping       map[uint32]*Manager
	PreparingPool *pool.Pool
	cipher        tunnel.Cipher
}

func NewCollector(cipher tunnel.Cipher) Collector {
	c := Collector{
		cipher:  cipher,
		mapping: make(map[uint32]*Manager),
	}
	c.PreparingPool = pool.NewPreparingPool(cipher, c.PushWorker)
	return c
}

func (collector *Collector) Serve(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		collector.PreparingPool.Manager.AddConn(conn)
	}
}

func (collector *Collector) PushWorker(clientID uint32, worker pool.Worker) {
	collector.lock.RLock()
	if man, ok := collector.mapping[clientID]; ok {
		collector.lock.RUnlock()
		man.RegisterWorker(worker)
	} else {
		p := pool.NewServerPool(clientID, collector.cipher)
		p.RegisterWorker(worker)
		man := NewConnectionManager(p)
		collector.lock.Lock()
		collector.mapping[clientID] = man
		collector.lock.Unlock()
	}
}
