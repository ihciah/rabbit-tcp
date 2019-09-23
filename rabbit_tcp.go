package rabbit_tcp

import (
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/pool"
	"net"
)

func NewRabbitTCP(config Config) *RabbitTCP {
	// TODO: Initialize:
	// 1. Use config to create pool
	connPool := pool.NewPool()
	return &RabbitTCP{
		Config: config,
		Pool:   connPool,
	}
}

type RabbitTCP struct {
	Config Config
	Pool   *pool.Pool
}

func (r *RabbitTCP) Dial(network, address string) (net.Conn, error) {
	return connection.NewRabbitTCPConn(r.Pool, network, address)
}
