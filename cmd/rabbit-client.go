package main

import (
	"bufio"
	"fmt"
	"github.com/ihciah/rabbit-tcp/client"
	"github.com/ihciah/rabbit-tcp/tunnel"
)

func startClient() {
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	cipher = nil
	c := client.NewClient(1, "127.0.0.1:9876", cipher)
	conn := c.Dial("www.baidu.com:80")
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}

func main() {
	startClient()
}
