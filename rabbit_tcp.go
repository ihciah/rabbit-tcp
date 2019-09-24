package rabbit_tcp

import (
	"github.com/ihciah/rabbit-tcp/config"
	"github.com/ihciah/rabbit-tcp/connection"
	"github.com/ihciah/rabbit-tcp/pool"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"math/rand"
	"net"
)

type RabbitTCPClient struct {
	Config      config.Config
	ConnManager *connection.Manager
}

func NewRabbitTCPClient(config config.Config) *RabbitTCPClient {
	// TODO: Use config to create pool
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	connPool := pool.NewClientPool(
		"127.0.0.1:12345",
		rand.Uint32(),
		cipher,
		[]uint{1},
	)
	connManager := connection.NewConnectionManager(connPool)
	return &RabbitTCPClient{
		Config:      config,
		ConnManager: connManager,
	}
}

func (r *RabbitTCPClient) Dial(address string) (net.Conn, error) {
	return connection.NewRabbitTCPConn(r.ConnManager, address)
}

type RabbitTCPServer struct {
	Config    config.Config
	Collector connection.Collector
}

func NewRabbitTCPServer(config config.Config) *RabbitTCPServer {
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	collector := connection.NewCollector(cipher)
	return &RabbitTCPServer{
		Config:    config,
		Collector: collector,
	}
}

func (r *RabbitTCPServer) Serve(address string) {
	r.Collector.Serve(address)
}
