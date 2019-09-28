package connection_pool

import (
	"github.com/ihciah/rabbit-tcp/connection"
	"net"
)

type ConnectionPool struct {
	connectionMapping map[uint32]connection.Connection
	manager           Manager
}

func (cp *ConnectionPool) AddConnection(conn net.Conn) {
	// 1. add to map
	// 2.
}

func (cp *ConnectionPool) RemoveConnection(conn net.Conn) {

}
