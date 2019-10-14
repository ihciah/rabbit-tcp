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

var Version = "No version information"

const (
	ClientMode = iota
	ServerMode
	DefaultPassword = "PASSWORD"
)

func parseFlags() (pass bool, mode int, password string, addr string, listen string, dest string, tunnelN int, verbose int) {
	var modeString string
	var printVersion bool
	flag.StringVar(&modeString, "mode", "c", "running mode(s or c)")
	flag.StringVar(&password, "password", DefaultPassword, "password")
	flag.StringVar(&addr, "rabbit-addr", ":443", "listen(server mode) or remote(client mode) address used by rabbit-tcp")
	flag.StringVar(&listen, "listen", "", "[Client Only] listen address, eg: 127.0.0.1:2333")
	flag.StringVar(&dest, "dest", "", "[Client Only] destination address, eg: shadowsocks server address")
	flag.IntVar(&tunnelN, "tunnelN", 4, "[Client Only] number of tunnels to use in rabbit-tcp")
	flag.IntVar(&verbose, "verbose", 2, "verbose level(0~5)")
	flag.BoolVar(&printVersion, "version", false, "show version")
	flag.Parse()

	pass = true

	// version
	if printVersion {
		log.Println("Rabbit TCP (https://github.com/ihciah/rabbit-tcp/)")
		log.Printf("Version: %s.\n", Version)
		pass = false
		return
	}

	// mode
	modeString = strings.ToLower(modeString)
	if modeString == "c" || modeString == "client" {
		mode = ClientMode
	} else if modeString == "s" || modeString == "server" {
		mode = ServerMode
	} else {
		log.Printf("Unsupported mode %s.\n", modeString)
		pass = false
		return
	}

	// password
	if password == "" {
		log.Println("Password must be specified.")
		pass = false
		return
	}
	if password == DefaultPassword {
		log.Println("Password must be changed instead of default password.")
		pass = false
		return
	}

	// listen, dest, tunnelN
	if mode == ClientMode {
		if listen == "" {
			log.Println("Listen address must be specified in client mode.")
			pass = false
		}
		if dest == "" {
			log.Println("Destination address must be specified in client mode.")
			pass = false
		}
		if tunnelN == 0 {
			log.Println("Tunnel number must be positive.")
			pass = false
		}
	}
	return
}

func main() {
	pass, mode, password, addr, listen, dest, tunnelN, verbose := parseFlags()
	if !pass {
		return
	}
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
