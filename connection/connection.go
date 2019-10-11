package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"go.uber.org/atomic"
	"net"
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
	closed           *atomic.Bool
	sendQueue        chan<- block.Block // Same as connectionPool
	recvQueue        chan block.Block
	orderedRecvQueue chan block.Block
	logger           *logger.Logger
}

func (bc *baseConnection) Stop() {
	bc.blockProcessor.removeFromPool()
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
	bc.logger.Debugf("Send connect to %s block.\n", address)
	blk := bc.blockProcessor.packConnect(address, bc.connectionID)
	bc.sendQueue <- blk
}

func (bc *baseConnection) SendDisconnect() {
	bc.logger.Debugln("Send disconnect block.")
	blk := bc.blockProcessor.packDisconnect(bc.connectionID)
	bc.sendQueue <- blk
	bc.Stop()
}

func (bc *baseConnection) sendData(data []byte) {
	bc.logger.Debugln("Send data block.")
	blocks := bc.blockProcessor.packData(data, bc.connectionID)
	for _, blk := range blocks {
		bc.sendQueue <- blk
	}
}
