package pool

import (
	"github.com/ihciah/rabbit-tcp/tunnel"
	"net"
	"time"
)

const (
	CHECK_INTERVAL_MS = 500
)

type Manager interface {
	AddConn(net.Conn)
	Daemon()
	Exit()
}

type basicManager struct {
	Pool     *Pool
	ExitChan chan bool
	Cipher   tunnel.Cipher
	Running  bool
}

func newBasicManager(pool *Pool, cipher tunnel.Cipher) basicManager {
	return basicManager{
		Pool:     pool,
		ExitChan: make(chan bool, 1),
		Cipher:   cipher,
		Running:  true,
	}
}

func (man *basicManager) AddConn(conn net.Conn) {
	tun := tunnel.NewTunnel(conn, man.Cipher)
	// TODO: tun need to send Reg block
	worker := NewBasicWorker(tun)
	man.Pool.RegisterWorker(worker)
}

func (man *basicManager) Exit() {
	if man.Running {
		man.ExitChan <- true
	}
}

type ClientManager struct {
	basicManager
	dest        string
	createUntil []time.Time
}

func NewClientManager(pool *Pool, dest string, cipher tunnel.Cipher, createWaitSec []uint) Manager {
	created := time.Now()
	createUntil := make([]time.Time, 0, len(createWaitSec))
	for _, sec := range createWaitSec {
		createUntil = append(createUntil, created.Add(time.Duration(sec)*time.Second))
	}
	return &ClientManager{
		basicManager: newBasicManager(pool, cipher),
		createUntil:  createUntil,
		dest:         dest,
	}
}

// Infinitely check pool size and fill the pool
func (man *ClientManager) Daemon() {
	for {
		select {
		case <-man.ExitChan:
			return
		case <-time.After(time.Millisecond * CHECK_INTERVAL_MS):
			man.keepPoolFull()
		}
	}
}

// Check if the pool size satisfies the config
// If not, create tunnel and add it to pool
func (man *ClientManager) keepPoolFull() {
	now := time.Now()
	i := 0
	for i = len(man.createUntil); i >= 0 && now.Before(man.createUntil[i]); i-- {
	}
	if i < 0 {
		return
	}
	tunnelToCreate := i - man.Pool.GetWorkerCount()

	for i = 0; i < tunnelToCreate; {
		err := man.createTunnel()
		if err != nil {
			time.Sleep(time.Second)
		} else {
			i++
		}
	}
}

// Create a tunnel and reg it; then add it to Workers
func (man *ClientManager) createTunnel() error {
	conn, err := net.Dial("tcp", man.dest)
	if err != nil {
		return err
	}
	man.AddConn(conn)
	return nil
}

type ServerManager struct {
	basicManager
}

func NewServerManager(pool *Pool, cipher tunnel.Cipher) Manager {
	return &ServerManager{
		basicManager: newBasicManager(pool, cipher),
	}
}

func (man *ServerManager) Daemon() {
	// TODO
}

type PreparingManager struct {
	basicManager
	pushFunc func(clientID uint32, worker Worker)
}

func NewPreparingManager(pool *Pool, cipher tunnel.Cipher, pushFunc func(clientID uint32, worker Worker)) Manager {
	return &PreparingManager{
		basicManager: newBasicManager(pool, cipher),
		pushFunc:     pushFunc,
	}
}

func (man *PreparingManager) Daemon() {
	// TODO: recv RegBlock and add connection to other pools
}
