package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"log"
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
	GetOrderedRecvQueue() chan block.Block
	GetRecvQueue() chan block.Block

	RecvBlock(block.Block)
	SendBlock(block.Block)

	SendData(data []byte)
	SendConnect(address string)
	SendDisconnect()

	Daemon(connection Connection)
	StopDaemon()
}

type BaseConnection struct {
	blockProcessor
	connectionID     uint32
	ok               bool
	sendQueue        chan<- block.Block // same as connectionPool
	recvQueue        chan block.Block
	orderedRecvQueue chan block.Block
	logger           *log.Logger
}

func (bc *BaseConnection) GetConnectionID() uint32 {
	return bc.connectionID
}

func (bc *BaseConnection) RecvBlock(blk block.Block) {
	bc.recvQueue <- blk
}

func (bc *BaseConnection) SendBlock(blk block.Block) {
	bc.sendQueue <- blk
}

func (bc *BaseConnection) GetRecvQueue() chan block.Block {
	return bc.recvQueue
}

func (bc *BaseConnection) GetOrderedRecvQueue() chan block.Block {
	return bc.orderedRecvQueue
}
