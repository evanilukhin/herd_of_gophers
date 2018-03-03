package main

import (
	"flag"
	"fmt"
	"github.com/evanilukhin/phochan"
	"time"
)

type Connection struct {
	Socket  *phochan.Socket
	Channel *phochan.Channel
}

func main() {
	firebusHost := flag.String("host", "0.0.0.0:4000", "FireBus address. Default: 0.0.0.0:4000")
	countClients := flag.Int("count", 1, "Count clients must be integer > 0. Default: 1")

	flag.Parse()

	serverAddress := "ws://" + *firebusHost + "/socket/websocket?token=undefined&vsn=2.0.0"
	mfunc := func(message []byte) { fmt.Println(string(message)) }
	var a []Connection
	for i := 0; i < *countClients; i++ {
		socket := phochan.NewSocket(serverAddress)
		channel := socket.Channel("test_room:lobby", mfunc)
		a = append(a, Connection{Socket: socket, Channel: channel})
	}
	for _, connection := range a {
		connection.Socket.Connect()
		connection.Channel.Join()
		connection.Channel.Start()
	}
	time.Sleep(time.Minute * 60)
}
