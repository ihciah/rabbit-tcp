package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"math/rand"
	"net"
	"time"
)

const (
	DataQueueSize = 24
	RecvQueueSize = 24
	BlockMaxSize  = 4096
)

type Connection struct {
	blockProcessor
	conn         net.Conn
	connectionID uint32
	dataQueue    chan block.Block
	dataBuffer   LoopByteBuffer
	sendQueue    chan<- block.Block // same as connectionPool
	recvQueue    chan block.Block
}

func NewConnection(conn net.Conn, sendQueue chan<- block.Block) Connection {
	connectionID := rand.Uint32()
	return NewConnectionWithID(conn, connectionID, sendQueue)
}

func NewConnectionWithID(conn net.Conn, connectionID uint32, sendQueue chan<- block.Block) Connection {
	c := Connection{
		conn:         conn,
		connectionID: connectionID,
		dataQueue:    make(chan block.Block, DataQueueSize),
		dataBuffer:   NewLoopBuffer(BlockMaxSize),
		sendQueue:    sendQueue,
		recvQueue:    make(chan block.Block, RecvQueueSize),
	}
	c.blockProcessor = newBlockProcessor(&c)
	return c
}

func (c *Connection) Read(b []byte) (n int, err error) {

}

func (c *Connection) Write(b []byte) (n int, err error) {
	// TODO: tag all blocks from b using WaitGroup
	// TODO: and wait all blocks sent
}

func (c *Connection) Close() error {

}

func (c *Connection) LocalAddr() net.Addr {

}

func (c *Connection) RemoteAddr() net.Addr {

}

func (c *Connection) SetDeadline(t time.Time) error {

}

func (c *Connection) SetReadDeadline(t time.Time) error {

}

func (c *Connection) SetWriteDeadline(t time.Time) error {

}
