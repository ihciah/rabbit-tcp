package rabbit_tcp

import (
	"github.com/ihciah/rabbit-tcp/config"
	"github.com/ihciah/rabbit-tcp/connection_old"
	"github.com/ihciah/rabbit-tcp/pool_old"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"math/rand"
	"net"
)

type RabbitTCPClient struct {
	Config      config.Config
	ConnManager *connection_old.Manager
}

func NewRabbitTCPClient(config config.Config) *RabbitTCPClient {
	// TODO: Use config to create pool
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	connPool := pool_old.NewClientPool(
		"127.0.0.1:12345",
		rand.Uint32(),
		cipher,
		[]uint{2},
	)
	connManager := connection_old.NewConnectionManager(connPool)
	return &RabbitTCPClient{
		Config:      config,
		ConnManager: connManager,
	}
}

func (r *RabbitTCPClient) Dial(address string) (net.Conn, error) {
	return connection_old.NewRabbitTCPConn(r.ConnManager, address)
}

type RabbitTCPServer struct {
	Config    config.Config
	Collector connection_old.Collector
}

func NewRabbitTCPServer(config config.Config) *RabbitTCPServer {
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	collector := connection_old.NewCollector(cipher)
	return &RabbitTCPServer{
		Config:    config,
		Collector: collector,
	}
}

func (r *RabbitTCPServer) Serve(address string) {
	r.Collector.Serve(address)
}
