package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"math/rand"
	"net"
	"time"
)

type InboundConnection struct {
	BaseConnection
	dataBuffer LoopByteBuffer
}

func NewInboundConnection(sendQueue chan<- block.Block) Connection {
	connectionID := rand.Uint32()
	return NewInboundConnectionWithID(connectionID, sendQueue)
}

func NewInboundConnectionWithID(connectionID uint32, sendQueue chan<- block.Block) Connection {
	c := InboundConnection{
		BaseConnection: BaseConnection{
			connectionID:     connectionID,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
		},
		dataBuffer: NewLoopBuffer(BlockMaxSize),
	}
	c.blockProcessor = newBlockProcessor(&c)
	return &c
}

func (c *InboundConnection) Read(b []byte) (n int, err error) {

}

func (c *InboundConnection) Write(b []byte) (n int, err error) {
	// TODO: tag all blocks from b using WaitGroup
	// TODO: and wait all blocks sent
}

func (c *InboundConnection) Close() error {

}

func (c *InboundConnection) LocalAddr() net.Addr {

}

func (c *InboundConnection) RemoteAddr() net.Addr {

}

func (c *InboundConnection) SetDeadline(t time.Time) error {

}

func (c *InboundConnection) SetReadDeadline(t time.Time) error {

}

func (c *InboundConnection) SetWriteDeadline(t time.Time) error {

}
