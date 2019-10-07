package tunnel_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	ErrorWaitSec        = 3
	EmptyPoolDestroySec = 60
)

type Manager interface {
	Notify(pool *TunnelPool)         // When TunnelPool size changed, Notify should be called
	DecreaseNotify(pool *TunnelPool) // When TunnelPool size decreased, DecreaseNotify should be called
}

type ClientManager struct {
	notifyLock sync.Mutex // Only one notify can run in the same time
	tunnelNum  int
	endpoint   string
	peerID     uint32
	cipher     tunnel.Cipher
	logger     *log.Logger
}

func NewClientManager(tunnelNum int, endpoint string, peerID uint32, cipher tunnel.Cipher) ClientManager {
	return ClientManager{
		tunnelNum: tunnelNum,
		endpoint:  endpoint,
		cipher:    cipher,
		peerID:    peerID,
		logger:    log.New(os.Stdout, "[ClientManager]", log.LstdFlags),
	}
}

// Keep tunnelPool size above tunnelNum
func (cm *ClientManager) DecreaseNotify(pool *TunnelPool) {
	cm.notifyLock.Lock()
	defer cm.notifyLock.Unlock()
	tunnelCount := len(pool.tunnelMapping)

	for tunnelToCreate := cm.tunnelNum - tunnelCount; tunnelToCreate > 0; tunnelToCreate-- {
		cm.logger.Printf("Need %d new tunnels now.\n", tunnelToCreate)
		conn, err := net.Dial("tcp", cm.endpoint)
		if err != nil {
			cm.logger.Printf("Error when dial to %s: %v.\n", cm.endpoint, err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		tun, err := NewActiveTunnel(conn, cm.cipher, cm.peerID)
		if err != nil {
			cm.logger.Printf("Error when create active tunnel: %v\n", err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		pool.AddTunnel(&tun)
		cm.logger.Printf("Successfully dialed to %s. TunnelToCreate: %d\n", cm.endpoint, tunnelToCreate)
	}
}

func (cm *ClientManager) Notify(pool *TunnelPool) {}

type ServerManager struct {
	notifyLock          sync.Mutex // Only one notify can run in the same time
	removePeerFunc      context.CancelFunc
	cancelCountDownFunc context.CancelFunc
	triggered           bool
	logger              *log.Logger
}

func NewServerManager(removePeerFunc context.CancelFunc) ServerManager {
	return ServerManager{
		logger:         log.New(os.Stdout, "[ServerManager]", log.LstdFlags),
		removePeerFunc: removePeerFunc,
	}
}

// If tunnelPool size is zero for more than EmptyPoolDestroySec, delete it
func (sm *ServerManager) Notify(pool *TunnelPool) {
	sm.notifyLock.Lock()
	defer sm.notifyLock.Unlock()
	tunnelCount := len(pool.tunnelMapping)

	if tunnelCount == 0 && !sm.triggered {
		var destroyAfterCtx context.Context
		destroyAfterCtx, sm.cancelCountDownFunc = context.WithCancel(context.Background())
		go func(*ServerManager) {
			select {
			case <-destroyAfterCtx.Done():
				sm.triggered = false
				sm.logger.Println("ServerManager notify canceled.")
				return
			case <-time.After(EmptyPoolDestroySec * time.Second):
				sm.logger.Println("ServerManager will be destroyed.")
				sm.removePeerFunc()
				return
			}
		}(sm)
	}

	if tunnelCount != 0 && sm.triggered {
		sm.cancelCountDownFunc()
	}
}

func (sm *ServerManager) DecreaseNotify(pool *TunnelPool) {}
