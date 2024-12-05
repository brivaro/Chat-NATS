package api

import (
	"fmt"        
	"log"        
	"time" 
    "context"  
	"chat/initializers"    

	"github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

func SubscribeToChannel(channel string) {
    subject := fmt.Sprintf(channel)
    _, err := initializers.Client.Conn.Subscribe(subject, func(msg *nats.Msg) {
        log.Printf(string(msg.Data))
    })
    if err != nil {
        log.Fatalf("Error subscribing channel %s: %v", channel, err)
    }
}

func PublishMessage(channel, user, message string) {
    subject := fmt.Sprintf("chat.%s", channel)
    fullMessage := fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), user, message)
    err := initializers.Client.Conn.Publish(subject, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error publishing mss: %v", err)
    }
}

func FetchRecentMessages(channel string) {
    lastHour := time.Now().Add(-1 * time.Hour)
    msgChan := make(chan *nats.Msg, 100)
    go func() {
        for {
            msg := <-msgChan
            fmt.Println(string(msg.Data))
        }
    }()

    go func() {
        initializers.JS.PurgeStream(streamName)
    }()
}