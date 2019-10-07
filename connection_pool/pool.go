package connection_pool

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/tunnel_pool"
	"log"
	"os"
	"sync"
)

const (
	SendQueueSize = 24
)

type ConnectionPool struct {
	connectionMapping map[uint32]connection.Connection
	mappingLock       sync.Mutex
	tunnelPool        *tunnel_pool.TunnelPool
	sendQueue         chan block.Block
	logger            *log.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewConnectionPool(pool *tunnel_pool.TunnelPool, backgroundCtx context.Context) ConnectionPool {
	ctx, cancel := context.WithCancel(backgroundCtx)
	cp := ConnectionPool{
		connectionMapping: make(map[uint32]connection.Connection),
		tunnelPool:        pool,
		sendQueue:         make(chan block.Block, SendQueueSize),
		logger:            log.New(os.Stdout, "[ConnectionPool]", log.LstdFlags),
		ctx:               ctx,
		cancel:            cancel,
	}
	cp.logger.Println("ConnectionPool created.")
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
	cp.logger.Printf("Connection %d created.\n", connectionID)
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
	cp.logger.Printf("Connection %d added to connection pool.\n", conn.GetConnectionID())
	cp.mappingLock.Lock()
	defer cp.mappingLock.Unlock()
	cp.connectionMapping[conn.GetConnectionID()] = conn
	go conn.OrderedRelay(conn)
}

func (cp *ConnectionPool) removeConnection(conn connection.Connection) {
	cp.logger.Printf("Connection %d removed from connection pool.\n", conn.GetConnectionID())
	cp.mappingLock.Lock()
	defer cp.mappingLock.Unlock()
	if _, ok := cp.connectionMapping[conn.GetConnectionID()]; ok {
		delete(cp.connectionMapping, conn.GetConnectionID())
	}
}

// Deliver blocks from tunnelPool channel to specified connections
func (cp *ConnectionPool) recvRelay() {
	cp.logger.Println("Recv Relay started.")
	for {
		select {
		case blk := <-cp.tunnelPool.GetRecvQueue():
			connID := blk.ConnectionID
			var conn connection.Connection
			var ok bool
			if conn, ok = cp.connectionMapping[connID]; !ok {
				conn = cp.NewPooledOutboundConnection(blk.ConnectionID)
				cp.logger.Println("Connection created and added to connectionPool.")
			}
			conn.RecvBlock(blk)
			cp.logger.Printf("Block %d(type: %d) put to connRecvQueue.\n", blk.BlockID, blk.Type)
		case <-cp.ctx.Done():
			cp.logger.Println("Recv Relay stoped.")
			return
		}
	}
}

// Deliver blocks from connPool's sendQueue to tunnelPool
// TODO: Maybe QOS can be implemented here
func (cp *ConnectionPool) sendRelay() {
	cp.logger.Println("Send Relay started.")
	for {
		select {
		case blk := <-cp.sendQueue:
			cp.tunnelPool.GetSendQueue() <- blk
			cp.logger.Printf("Block %d(type: %d) put to connSendQueue.\n", blk.BlockID, blk.Type)
		case <-cp.ctx.Done():
			cp.logger.Printf("Send Relay stoped.\n")
			return
		}
	}
}

func (cp *ConnectionPool) stopRelay() {
	cp.logger.Println("Stop all ConnectionPool Relay.")
	cp.cancel()
}
