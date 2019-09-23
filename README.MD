# Rabbit TCP
A multi-connection TCP

## Usage
Create a instance with `rabbitTCP := RabbitTCP.New()`

Use `rabbitTCP.Dial` as `net.Dial`

## Example

```go
package main
import "github.com/ihciah/rabbit-tcp"

config := RabbitTCP.LoadConfigFromFile("config.json")
rabbitTCP := RabbitTCP.NewRabbitTCP(config)
conn, err := rabbitTCP.Dial("tcp", "golang.org:80")
if err != nil {
	// handle error
}
fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
status, err := bufio.NewReader(conn).ReadString('\n')
```