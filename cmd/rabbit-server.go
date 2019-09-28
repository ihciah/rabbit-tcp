package main

import (
	rabbit_tcp "github.com/ihciah/rabbit-tcp"
	"github.com/ihciah/rabbit-tcp/config"
)

func startServer() {
	conf := config.LoadConfigFromFile("config.json")
	rabbitTCP := rabbit_tcp.NewRabbitTCPServer(conf)
	rabbitTCP.Serve("127.0.0.1:12345")
}

func main() {
	startServer()
}
