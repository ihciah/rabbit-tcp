package pool

import (
	"github.com/ihciah/rabbit-tcp/tunnel"
	"net"
	"sync"
)

type Collector struct {
	lock          sync.RWMutex
	mapping       map[uint32]*Pool
	PreparingPool *Pool
	cipher        tunnel.Cipher
}

func NewCollector(cipher tunnel.Cipher) Collector {
	c := Collector{
		cipher:  cipher,
		mapping: make(map[uint32]*Pool),
	}
	c.PreparingPool = NewPreparingPool(cipher, c.PushWorker)
	return c
}

func (collector *Collector) Listen(network, address string) {
	listener, err := net.Listen(network, address)
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

func (collector *Collector) PushWorker(clientID uint32, worker Worker) {
	collector.lock.RLock()
	if p, ok := collector.mapping[clientID]; ok {
		collector.lock.RUnlock()
		p.RegisterWorker(worker)
	} else {
		p = NewServerPool(clientID, collector.cipher)
		p.RegisterWorker(worker)
		collector.lock.Lock()
		collector.mapping[clientID] = p
		collector.lock.Unlock()
	}
}
