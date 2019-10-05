package main

import (
	"fmt"
	"github.com/ihciah/rabbit-tcp/client"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"io/ioutil"
	"time"
)

func startClient() {
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, "password")
	cipher = nil
	c := client.NewClient(6, "127.0.0.1:9876", cipher)
	time.Sleep(3 * time.Second)
	conn := c.Dial("www.baidu.com:80")
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	result, err := ioutil.ReadAll(conn)
	fmt.Printf("Error: %v, Length: %d\n", err, len(result))
	fmt.Println(string(result))
}

func main() {
	startClient()
}
