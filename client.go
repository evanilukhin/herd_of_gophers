package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/evanilukhin/phochan"
	"time"
)

// Connection contains link to socket with channel for easy manipulation clients
// connected to same topics
type Connection struct {
	Socket  *phochan.Socket
	Channel *phochan.Channel
}

//IncomingPayload message from firebus
type IncomingPayload struct {
	Item      int    `json:"item"`
	Firebus   string `json:"firebus"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
}

//OutcomingPayload  IncomingPayload message after adding current application timestamp
type OutcomingPayload struct {
	Item          int    `json:"item"`
	Firebus       string `json:"firebus"`
	UUID          string `json:"uuid"`
	HerdOfGophers string `json:"herd_of_gophers"`
	CreatedAt     string `json:"created_at"`
}

var topic = "herd_of_gophers"
var collector = make(chan []byte)

func main() {
	firebusHost := flag.String("host", "0.0.0.0:4000", "FireBus address. Default: 0.0.0.0:4000")
	countClients := flag.Int("count", 1, "Count clients must be integer > 0. Default: 1")
	kafkaBrokerHost := flag.String("kafka", "0.0.0.0:9092", "Kafka broker address. Default: 0.0.0.0:9092")

	flag.Parse()

	serverAddress := "ws://" + *firebusHost + "/socket/websocket?token=undefined&vsn=2.0.0"

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": *kafkaBrokerHost})
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	var a []Connection
	for i := 0; i < *countClients; i++ {
		socket := phochan.NewSocket(serverAddress)
		channel := socket.Channel("test_room:lobby", transformAndPrint)
		a = append(a, Connection{Socket: socket, Channel: channel})
	}
	for _, connection := range a {
		connection.Socket.Connect()
		connection.Channel.Join()
		connection.Channel.Start()
	}

	go func() {
		for word := range collector {
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          word,
			}, nil)
		}
	}()
	time.Sleep(time.Minute * 60)
}

func transformAndPrint(message phochan.PhoenixMessage) {
	t := time.Now().UTC()
	var payload IncomingPayload
	json.Unmarshal(message.Payload, &payload)
	outcomingPayload := OutcomingPayload{
		Item:          payload.Item,
		UUID:          payload.UUID,
		Firebus:       payload.Firebus,
		CreatedAt:     payload.CreatedAt,
		HerdOfGophers: t.Format("15:04:05.999999"),
	}
	marshalledOutcomingPayload, _ := json.Marshal(outcomingPayload)
	collector <- marshalledOutcomingPayload
}
