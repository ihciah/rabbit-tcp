package main

import (
	"bufio"
	"fmt"
	rabbit_tcp "github.com/ihciah/rabbit-tcp"
	"github.com/ihciah/rabbit-tcp/config"
	"time"
)

func startClient() {
	conf := config.LoadConfigFromFile("config.json")
	rabbitTCP := rabbit_tcp.NewRabbitTCPClient(conf)
	time.Sleep(5 * time.Second)
	conn, err := rabbitTCP.Dial("www.baidu.com:80")
	fmt.Println("err:", err)
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("err:", err)
	fmt.Println(status)
}

func main() {
	startClient()
}
