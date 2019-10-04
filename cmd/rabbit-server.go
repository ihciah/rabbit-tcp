package main

import (
	"github.com/ihciah/rabbit-tcp/server"
	"github.com/ihciah/rabbit-tcp/tunnel"
)

func startServer() {
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	cipher = nil
	server := server.NewServer(cipher)
	server.Serve("127.0.0.1:9876")
}

func main() {
	startServer()
}
