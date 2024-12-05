package initializers

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	Conn        *nats.Conn
}

var Client NATSClient

func ConnectToNats(serverURL string) {
    nc, err := nats.Connect(serverURL) // nats://localhost:4222 este es el default
    if err != nil {
        log.Fatalf("Failed to connect to NATS server: %v", err)
    }
    Client.Conn = nc
	log.Printf("Connected to NATS")
}