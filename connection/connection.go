package connection

import (
	"math/rand"
	"net"
	"time"
)

const CONNECTION_BUFFER = 65500 // Must be bigger or equal than block size
const DATA_QUEUE_SIZE = 24

// Fill in clientID reading from pool, Generate ConnID, send Connect Block
func NewRabbitTCPConn(connManager *Manager, address string) (*RabbitTCPConn, error) {
	c := &RabbitTCPConn{
		clientID:     connManager.ClientID,
		connectionID: rand.Uint32(),
		connManager:  connManager,
		recvQueue:    make(chan []byte, DATA_QUEUE_SIZE),
		recvBuffer:   NewLoopBuffer(CONNECTION_BUFFER),
		sendBlockID:  1,
		recvBlockID:  1,
		address:      address,
	}
	connManager.Connect(c, address)
	return c, nil
}

type RabbitTCPConn struct {
	clientID     uint32
	connectionID uint32
	connManager  *Manager
	recvQueue    chan []byte
	recvBuffer   LoopByteBuffer
	sendBlockID  uint32
	recvBlockID  uint32
	address      string
}

func (conn *RabbitTCPConn) Read(b []byte) (r int, err error) {
	left := len(b)
	copied := 0

	// Read data from recvBuffer first
	if !conn.recvBuffer.Empty() {
		n := conn.recvBuffer.Read(b)
		copied += n
		left -= n
	}

	// Buffer is empty now
	for left > 0 {
		select {
		// Can get block from recvQueue
		case blockData := <-conn.recvQueue:
			var n int
			// Copy to dst directly
			if len(blockData) <= left {
				n = copy(b[copied:], blockData)
			} else {
				conn.recvBuffer.OverWrite(blockData)
				n = conn.recvBuffer.Read(b[copied:])
			}
			copied += n
			left -= n
		// Queue is empty now, should return read data directly
		default:
			return copied, nil
		}
	}
	return copied, nil
}

func (conn *RabbitTCPConn) Write(b []byte) (n int, err error) {
	// Send to connManager
	return conn.connManager.SendBytes(conn, b)
}

func (conn *RabbitTCPConn) Close() error {

}

func (conn *RabbitTCPConn) LocalAddr() net.Addr {

}

func (conn *RabbitTCPConn) RemoteAddr() net.Addr {

}

func (conn *RabbitTCPConn) SetDeadline(t time.Time) error {

}

func (conn *RabbitTCPConn) SetReadDeadline(t time.Time) error {

}

func (conn *RabbitTCPConn) SetWriteDeadline(t time.Time) error {

}
