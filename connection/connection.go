package connection

import (
	"github.com/ihciah/rabbit-tcp/pool"
	"net"
	"time"
)

func NewRabbitTCPConn(pool *pool.Pool, network, address string) (*RabbitTCPConn, error){

}

type RabbitTCPConn struct {

}

func (conn *RabbitTCPConn) Read(b []byte) (n int, err error){

}

func (conn *RabbitTCPConn) Write(b []byte) (n int, err error){

}

func (conn *RabbitTCPConn) Close() error{

}

func (conn *RabbitTCPConn) LocalAddr() net.Addr{

}

func (conn *RabbitTCPConn) RemoteAddr() net.Addr{

}

func (conn *RabbitTCPConn) SetDeadline(t time.Time) error{

}

func (conn *RabbitTCPConn) SetReadDeadline(t time.Time) error{

}

func (conn *RabbitTCPConn) SetWriteDeadline(t time.Time) error{

}