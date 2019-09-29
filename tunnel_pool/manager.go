package tunnel_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"log"
	"net"
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
	cipher     tunnel.Cipher
	logger     log.Logger
}

// Keep tunnelPool size above tunnelNum
func (cm *ClientManager) DecreaseNotify(pool *TunnelPool) {
	cm.notifyLock.Lock()
	defer cm.notifyLock.Unlock()
	tunnelCount := pool.GetTunnelCount()

	for tunnelToCreate := cm.tunnelNum - tunnelCount; tunnelToCreate > 0; tunnelToCreate-- {
		conn, err := net.Dial("tcp", cm.endpoint)
		if err != nil {
			cm.logger.Printf("Can not dial to %s. Cause: %s\n", cm.endpoint, err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		tun, err := NewTunnel(conn, cm.cipher)
		if err != nil {
			cm.logger.Printf("Can not create tunnel. Cause: %s\n", err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		pool.AddTunnel(tun)
		cm.logger.Printf("Successfully dialed to %s. TunnelToCreate: %d\n", cm.endpoint, tunnelToCreate)
	}
}

func (cm *ClientManager) Notify(pool *TunnelPool) {}

type ServerManager struct {
	notifyLock          sync.Mutex // Only one notify can run in the same time
	destroyMeFunc       context.CancelFunc
	cancelCountDownFunc context.CancelFunc
	triggered           bool
	cipher              tunnel.Cipher
	logger              log.Logger
}

// If tunnelPool size is zero for more than EmptyPoolDestroySec, delete it
func (sm *ServerManager) Notify(pool *TunnelPool) {
	sm.notifyLock.Lock()
	defer sm.notifyLock.Unlock()
	tunnelCount := pool.GetTunnelCount()

	if tunnelCount == 0 && !sm.triggered {
		var destroyAfterCtx context.Context
		destroyAfterCtx, sm.cancelCountDownFunc = context.WithCancel(context.Background())
		go func(*ServerManager) {
			select {
			case <-destroyAfterCtx.Done():
				sm.triggered = false
				return
			case <-time.After(EmptyPoolDestroySec * time.Second):
				sm.destroyMeFunc()
				return
			}
		}(sm)
	}

	if tunnelCount != 0 && sm.triggered {
		sm.cancelCountDownFunc()
	}
}

func (sm *ServerManager) DecreaseNotify(pool *TunnelPool) {}
