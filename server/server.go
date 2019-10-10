package server

import (
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"log"
	"net"
	"os"
)

type Server struct {
	peerGroup peer.PeerGroup
	logger    *log.Logger
}

func NewServer(cipher tunnel.Cipher) Server {
	return Server{
		peerGroup: peer.NewPeerGroup(cipher),
		logger:    log.New(os.Stdout, "[Server]", log.LstdFlags),
	}
}

func (s *Server) Serve(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Printf("Error when accept connection: %v.\n", err)
		}
		err = s.peerGroup.AddTunnelFromConn(conn)
		if err != nil {
			s.logger.Printf("Error when add tunnel to tunnel pool: %v.\n", err)
		}
	}
}
