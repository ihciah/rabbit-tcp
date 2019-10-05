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
	getOrderedRecvQueue() chan block.Block
	getRecvQueue() chan block.Block

	RecvBlock(block.Block)
	sendBlock(block.Block)

	sendData(data []byte)
	SendConnect(address string)
	SendDisconnect()

	Daemon(connection Connection)
	StopDaemon()
}

type baseConnection struct {
	blockProcessor
	connectionID     uint32
	ok               bool
	sendQueue        chan<- block.Block // same as connectionPool
	recvQueue        chan block.Block
	orderedRecvQueue chan block.Block
	logger           *log.Logger
}

func (bc *baseConnection) GetConnectionID() uint32 {
	return bc.connectionID
}

func (bc *baseConnection) RecvBlock(blk block.Block) {
	bc.recvQueue <- blk
}

func (bc *baseConnection) sendBlock(blk block.Block) {
	bc.sendQueue <- blk
}

func (bc *baseConnection) getRecvQueue() chan block.Block {
	return bc.recvQueue
}

func (bc *baseConnection) getOrderedRecvQueue() chan block.Block {
	return bc.orderedRecvQueue
}

func (bc *baseConnection) sendData(data []byte) {
	bc.logger.Printf("Send data block.\n")
	blocks := bc.packData(data, bc.connectionID)
	for _, blk := range blocks {
		bc.sendBlock(blk)
	}
}

func (bc *baseConnection) SendConnect(address string) {
	bc.logger.Printf("Send connect to %s block.\n", address)
	blk := bc.packConnect(address, bc.connectionID)
	bc.sendBlock(blk)
}

func (bc *baseConnection) SendDisconnect() {
	bc.logger.Printf("Send disconnect block.\n")
	blk := bc.packDisconnect(bc.connectionID)
	bc.sendBlock(blk)
}
