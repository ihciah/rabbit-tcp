package tunnel_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

type Manager interface {
	Notify(pool *TunnelPool)         // When TunnelPool size changed, Notify should be called
	DecreaseNotify(pool *TunnelPool) // When TunnelPool size decreased, DecreaseNotify should be called
}

type ClientManager struct {
	decreaseNotifyLock sync.Mutex // Only one decrease notify can run at the same time
	tunnelNum          int
	endpoint           string
	peerID             uint32
	cipher             tunnel.Cipher
	logger             *logger.Logger
}

func NewClientManager(tunnelNum int, endpoint string, peerID uint32, cipher tunnel.Cipher) ClientManager {
	return ClientManager{
		tunnelNum: tunnelNum,
		endpoint:  endpoint,
		cipher:    cipher,
		peerID:    peerID,
		logger:    logger.NewLogger("[ClientManager]"),
	}
}

// Keep tunnelPool size above tunnelNum
func (cm *ClientManager) DecreaseNotify(pool *TunnelPool) {
	cm.decreaseNotifyLock.Lock()
	defer cm.decreaseNotifyLock.Unlock()
	tunnelCount := len(pool.tunnelMapping)

	for tunnelToCreate := cm.tunnelNum - tunnelCount; tunnelToCreate > 0; {
		select {
		case <-pool.ctx.Done():
			// Have to return if pool cancel is called.
			return
		default:
		}

		cm.logger.Infof("Need %d new tunnels now.\n", tunnelToCreate)
		conn, err := net.Dial("tcp", cm.endpoint)
		if err != nil {
			cm.logger.Errorf("Error when dial to %s: %v.\n", cm.endpoint, err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		tun, err := NewActiveTunnel(conn, cm.cipher, cm.peerID)
		if err != nil {
			cm.logger.Errorf("Error when create active tunnel: %v\n", err)
			time.Sleep(ErrorWaitSec * time.Second)
			continue
		}
		pool.AddTunnel(&tun)
		tunnelToCreate--
		cm.logger.Infof("Successfully dialed to %s. TunnelToCreate: %d\n", cm.endpoint, tunnelToCreate)
	}
}

func (cm *ClientManager) Notify(pool *TunnelPool) {}

type ServerManager struct {
	notifyLock          sync.Mutex // Only one notify can run in the same time
	removePeerFunc      context.CancelFunc
	cancelCountDownFunc context.CancelFunc
	triggered           atomic.Bool
	logger              *logger.Logger
}

func NewServerManager(removePeerFunc context.CancelFunc) ServerManager {
	return ServerManager{
		logger:         logger.NewLogger("[ServerManager]"),
		removePeerFunc: removePeerFunc,
	}
}

// If tunnelPool size is zero for more than EmptyPoolDestroySec, delete it
func (sm *ServerManager) Notify(pool *TunnelPool) {
	tunnelCount := len(pool.tunnelMapping)

	if tunnelCount == 0 && sm.triggered.CAS(false, true) {
		var destroyAfterCtx context.Context
		destroyAfterCtx, sm.cancelCountDownFunc = context.WithCancel(context.Background())
		go func(*ServerManager) {
			select {
			case <-destroyAfterCtx.Done():
				sm.logger.Debugln("ServerManager notify canceled.")
			case <-time.After(EmptyPoolDestroySec * time.Second):
				sm.logger.Infoln("ServerManager will be destroyed.")
				sm.removePeerFunc()
			}
		}(sm)
	}

	if tunnelCount != 0 && sm.triggered.CAS(true, false) {
		sm.cancelCountDownFunc()
	}
}

func (sm *ServerManager) DecreaseNotify(pool *TunnelPool) {}
