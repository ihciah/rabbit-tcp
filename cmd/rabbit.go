package main

import (
	"flag"
	"github.com/ihciah/rabbit-tcp/client"
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/server"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"log"
	"strings"
)

const (
	ClientMode = iota
	ServerMode
	DefaultPassword = "PASSWORD"
)

func parseFlags() (int, string, string, string, string, int, int) {
	var mode int
	var modeString string
	var password string
	var addr string
	var listen string
	var dest string
	var tunnelN int
	var verbose int
	flag.StringVar(&modeString, "mode", "c", "running mode(s or c)")
	flag.StringVar(&password, "password", DefaultPassword, "password")
	flag.StringVar(&addr, "rabbit-addr", ":443", "listen(server mode) or remote(client mode) address used by rabbit-tcp")
	flag.StringVar(&listen, "listen", "", "[Client Only] listen address, eg: 127.0.0.1:2333")
	flag.StringVar(&dest, "dest", "", "[Client Only] destination address, eg: shadowsocks server address")
	flag.IntVar(&tunnelN, "tunnelN", 6, "[Client Only] number of tunnels to use in rabbit-tcp")
	flag.IntVar(&verbose, "verbose", 2, "Verbose level(0~5)")
	flag.Parse()

	modeString = strings.ToLower(modeString)
	if modeString == "c" || modeString == "client" {
		mode = ClientMode
	} else if modeString == "s" || modeString == "server" {
		mode = ServerMode
	} else {
		log.Panicf("Unsupported mode %s.\n", modeString)
	}

	if password == DefaultPassword {
		log.Panicln("Password must be changed instead of default password.")
	}

	if mode == ClientMode {
		if listen == "" {
			log.Panicln("Listen address must be specified in client mode.")
		}
		if dest == "" {
			log.Panicln("Destination address must be specified in client mode.")
		}
		if tunnelN == 0 {
			log.Panicln("Tunnel number must be positive.")
		}
	}
	return mode, password, addr, listen, dest, tunnelN, verbose
}

func main() {
	mode, password, addr, listen, dest, tunnelN, verbose := parseFlags()
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, password)
	logger.LEVEL = verbose
	if mode == ClientMode {
		c := client.NewClient(tunnelN, addr, cipher)
		c.ServeForward(listen, dest)
	} else {
		s := server.NewServer(cipher)
		s.Serve(addr)
	}
}
