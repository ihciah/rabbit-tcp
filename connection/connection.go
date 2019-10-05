package connection

import (
	"context"
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

	SendConnect(address string)
	SendDisconnect()

	OrderedRelay(connection Connection) // Run orderedRelay infinitely
	Stop()                              // Stop all related relay and remove itself from connectionPool
}

type baseConnection struct {
	blockProcessor   blockProcessor
	connectionID     uint32
	ok               bool
	sendQueue        chan<- block.Block // Same as connectionPool
	recvQueue        chan block.Block
	orderedRecvQueue chan block.Block
	removeFromPool   context.CancelFunc
	logger           *log.Logger
}

func (bc *baseConnection) Stop() {
	if bc.removeFromPool != nil {
		bc.removeFromPool()
	}
}

func (bc *baseConnection) OrderedRelay(connection Connection) {
	bc.blockProcessor.OrderedRelay(connection)
}

func (bc *baseConnection) GetConnectionID() uint32 {
	return bc.connectionID
}

func (bc *baseConnection) getRecvQueue() chan block.Block {
	return bc.recvQueue
}

func (bc *baseConnection) getOrderedRecvQueue() chan block.Block {
	return bc.orderedRecvQueue
}

func (bc *baseConnection) RecvBlock(blk block.Block) {
	bc.recvQueue <- blk
}

func (bc *baseConnection) SendConnect(address string) {
	bc.logger.Printf("Send connect to %s block.\n", address)
	blk := bc.blockProcessor.packConnect(address, bc.connectionID)
	bc.sendQueue <- blk
}

func (bc *baseConnection) SendDisconnect() {
	bc.logger.Printf("Send disconnect block.\n")
	blk := bc.blockProcessor.packDisconnect(bc.connectionID)
	bc.sendQueue <- blk
}

func (bc *baseConnection) sendData(data []byte) {
	bc.logger.Printf("Send data block.\n")
	blocks := bc.blockProcessor.packData(data, bc.connectionID)
	for _, blk := range blocks {
		bc.sendQueue <- blk
	}
}
