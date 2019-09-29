package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"net"
)

const (
	OrderedRecvQueueSize = 24
	RecvQueueSize        = 24
	BlockMaxSize         = 4096
)

type Connection interface {
	net.Conn
	GetConnectionID() uint32
	GetRecvQueue() chan block.Block
	GetSendQueue() chan<- block.Block
	GetOrderedRecvQueue() chan block.Block
}

type BaseConnection struct {
	blockProcessor
	connectionID     uint32
	sendQueue        chan<- block.Block // same as connectionPool
	recvQueue        chan block.Block
	orderedRecvQueue chan block.Block
}

func (bc *BaseConnection) GetConnectionID() uint32 {
	return bc.connectionID
}

func (bc *BaseConnection) GetRecvQueue() chan block.Block {
	return bc.recvQueue
}

func (bc *BaseConnection) GetSendQueue() chan<- block.Block {
	return bc.sendQueue
}

func (bc *BaseConnection) GetOrderedRecvQueue() chan block.Block {
	return bc.orderedRecvQueue
}
