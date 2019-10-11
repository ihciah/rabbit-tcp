package connection_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"sync"
)

const (
	SendQueueSize = 48 // SendQueue channel cap
)

type ConnectionPool struct {
	connectionMapping map[uint32]connection.Connection
	mappingLock       sync.Mutex
	tunnelPool        *tunnel_pool.TunnelPool
	sendQueue         chan block.Block
	logger            *logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewConnectionPool(pool *tunnel_pool.TunnelPool, backgroundCtx context.Context) ConnectionPool {
	ctx, cancel := context.WithCancel(backgroundCtx)
	cp := ConnectionPool{
		connectionMapping: make(map[uint32]connection.Connection),
		tunnelPool:        pool,
		sendQueue:         make(chan block.Block, SendQueueSize),
		logger:            logger.NewLogger("[ConnectionPool]"),
		ctx:               ctx,
		cancel:            cancel,
	}
	cp.logger.Infoln("Connection Pool created.")
	go cp.sendRelay()
	go cp.recvRelay()
	return cp
}

// Create InboundConnection, and it to ConnectionPool and return
func (cp *ConnectionPool) NewPooledInboundConnection() connection.Connection {
	connCtx, removeConnFromPool := context.WithCancel(cp.ctx)
	c := connection.NewInboundConnection(cp.sendQueue, connCtx, removeConnFromPool)
	cp.addConnection(c)
	go func() {
		<-connCtx.Done()
		cp.removeConnection(c)
	}()
	return c
}

// Create OutboundConnection, and it to ConnectionPool and return
func (cp *ConnectionPool) NewPooledOutboundConnection(connectionID uint32) connection.Connection {
	connCtx, removeConnFromPool := context.WithCancel(cp.ctx)
	c := connection.NewOutboundConnection(connectionID, cp.sendQueue, connCtx, removeConnFromPool)
	cp.addConnection(c)
	go func() {
		<-connCtx.Done()
		cp.removeConnection(c)
	}()
	return c
}

func (cp *ConnectionPool) addConnection(conn connection.Connection) {
	cp.logger.Infof("Connection %d added to connection pool.\n", conn.GetConnectionID())
	cp.mappingLock.Lock()
	defer cp.mappingLock.Unlock()
	cp.connectionMapping[conn.GetConnectionID()] = conn
	go conn.OrderedRelay(conn)
}

func (cp *ConnectionPool) removeConnection(conn connection.Connection) {
	cp.logger.Infof("Connection %d removed from connection pool.\n", conn.GetConnectionID())
	cp.mappingLock.Lock()
	defer cp.mappingLock.Unlock()
	if _, ok := cp.connectionMapping[conn.GetConnectionID()]; ok {
		delete(cp.connectionMapping, conn.GetConnectionID())
	}
}

// Deliver blocks from tunnelPool channel to specified connections
func (cp *ConnectionPool) recvRelay() {
	cp.logger.Infoln("Recv Relay started.")
	for {
		select {
		case blk := <-cp.tunnelPool.GetRecvQueue():
			connID := blk.ConnectionID
			var conn connection.Connection
			var ok bool
			if conn, ok = cp.connectionMapping[connID]; !ok {
				conn = cp.NewPooledOutboundConnection(blk.ConnectionID)
				cp.logger.Infoln("Connection created and added to connectionPool.")
			}
			conn.RecvBlock(blk)
			cp.logger.Debugf("Block %d(type: %d) put to connRecvQueue.\n", blk.BlockID, blk.Type)
		case <-cp.ctx.Done():
			cp.logger.Infoln("Recv Relay stopped.")
			return
		}
	}
}

// Deliver blocks from connPool's sendQueue to tunnelPool
// TODO: Maybe QOS can be implemented here
func (cp *ConnectionPool) sendRelay() {
	cp.logger.Infoln("Send Relay started.")
	for {
		select {
		case blk := <-cp.sendQueue:
			cp.tunnelPool.GetSendQueue() <- blk
			cp.logger.Debugf("Block %d(type: %d) put to connSendQueue.\n", blk.BlockID, blk.Type)
		case <-cp.ctx.Done():
			cp.logger.Infoln("Send Relay stopped.")
			return
		}
	}
}

func (cp *ConnectionPool) stopRelay() {
	cp.logger.Infoln("Stop all ConnectionPool Relay.")
	cp.cancel()
}
