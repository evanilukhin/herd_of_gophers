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

var firebusHost string

func main() {
	flag.StringVar(&firebusHost, "host", "0.0.0.0:4000", "FireBus address. Default: 0.0.0.0:4000")
	serverAddress := "ws://" + firebusHost + "/socket/websocket?token=undefined&vsn=2.0.0"
	mfunc := func(message string) { fmt.Println(message) }
	var a []Connection
	for i := 0; i < 100; i++ {
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
