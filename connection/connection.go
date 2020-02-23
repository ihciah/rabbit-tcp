package connection

import (
	"net"

	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"go.uber.org/atomic"
)

type HalfOpenConn interface {
	net.Conn
	CloseRead() error
	CloseWrite() error
}

type CloseWrite interface {
	CloseWrite() error
}

type CloseRead interface {
	CloseRead() error
}

type Connection interface {
	HalfOpenConn
	GetConnectionID() uint32
	getOrderedRecvQueue() chan block.Block
	getRecvQueue() chan block.Block

	RecvBlock(block.Block)

	SendConnect(address string)
	SendDisconnect(uint8)

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
	bc.logger.Debugf("connection stop\n")
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

func (bc *baseConnection) SendDisconnect(shutdownType uint8) {
	bc.logger.Debugf("Send disconnect block: %v\n", shutdownType)
	blk := bc.blockProcessor.packDisconnect(bc.connectionID, shutdownType)
	bc.sendQueue <- blk
	if shutdownType == block.ShutdownBoth {
		bc.Stop()
	}
}

func (bc *baseConnection) sendData(data []byte) {
	bc.logger.Debugln("Send data block.")
	blocks := bc.blockProcessor.packData(data, bc.connectionID)
	for _, blk := range blocks {
		bc.sendQueue <- blk
	}
}
