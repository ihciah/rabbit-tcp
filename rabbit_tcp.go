package rabbit_tcp

import (
	"github.com/ihciah/rabbit-tcp/config"
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/pool"
	"net"
)

func NewRabbitTCPClient(config config.Config) *RabbitTCP {
	// TODO: Use config to create pool
	connPool := pool.NewClientPool()
	return &RabbitTCP{
		Config: config,
		Pool:   connPool,
	}
}

type RabbitTCP struct {
	Config config.Config
	Pool   *pool.Pool
}

func (r *RabbitTCP) Dial(network, address string) (net.Conn, error) {
	return connection.NewRabbitTCPConn(r.Pool, network, address)
}
