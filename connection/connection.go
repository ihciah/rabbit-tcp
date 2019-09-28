package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"net"
)

type Connection struct {
	net.Conn
	connectionID   uint32
	blockProcessor blockProcessor
	dataQueue      chan block.Block
	dataBuffer     LoopByteBuffer
}

func (c *Connection) RecvBlock(block block.Block) {
	c.blockProcessor.blockBuffer <- block
}
