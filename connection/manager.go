package connection

import (
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/pool"
)

type Manager struct {
	*pool.Pool
	connJoiner map[uint32]*BlockJoiner
}

func NewConnectionManager(pool *pool.Pool) *Manager {
	return &Manager{
		Pool:       pool,
		connJoiner: make(map[uint32]*BlockJoiner),
	}
}

// pack data to block and send
func (m *Manager) SendBytes(conn *RabbitTCPConn, data []byte) (n int, err error) {
	blocks := block.NewDataBlocks(conn.connectionID, conn.sendBlockID, data)
	conn.sendBlockID += uint32(len(blocks))
	for _, b := range blocks {
		m.Pool.SendBlock(b)
	}
	return len(data), nil
}

// pack connect data to block and send
// create a joinMapping for the connection
// add connectionID->{conn, joinMapping} to connMapping
func (m *Manager) Connect(conn *RabbitTCPConn, address string) {
	connBlock := block.NewConnectBlock(conn.connectionID, conn.address)
	blockJoiner := NewBlockJoiner(conn)
	m.connJoiner[conn.connectionID] = blockJoiner
	m.Pool.SendBlock(connBlock)
}

func (m *Manager) Disconnect(conn *RabbitTCPConn) {
	if _, ok := m.connJoiner[conn.connectionID]; ok {
		delete(m.connJoiner, conn.connectionID)
	}
	disconnectBlock := block.NewDisconnectBlock(conn.connectionID)
	m.Pool.SendBlock(disconnectBlock)
}

func (m *Manager) recv() {
	blk := m.Pool.RecvBlock()
	switch blk.Type {
	case block.BLOCK_TYPE_CONNECT:
		// TODO
	case block.BLOCK_TYPE_DISCONNECT:
		// TODO
	case block.BLOCK_TYPE_REG:
		// TODO
	case block.BLOCK_TYPE_DATA:
		// receive a block from pool's recvBuffer
		// loop up mapping to find a {connection, joinMapping}
		// put block to joinMapping
		// if the data is ready to send, send it and remove from joinMapping
		connID := blk.ConnectionID
		if joiner, ok := m.connJoiner[connID]; ok {
			joiner.AddBlock(&blk)
			if data := joiner.GetBlock(); data != nil {
				joiner.conn.recvQueue <- data
			}
		}
	}
}

func (m *Manager) RecvDaemon() {
	// run a daemon to recv
	go func(m *Manager) {
		for {
			m.recv()
		}
	}(m)
}
