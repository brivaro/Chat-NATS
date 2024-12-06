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
    subject := fmt.Sprintf(channel)
    fullMessage := fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), user, message)
    err := initializers.Client.Conn.Publish(subject, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error publishing mss: %v", err)
    }
}

func FetchRecentMessages(channel string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf(channel) //("chat.%s", channel)

    startTime := time.Now().Add(-1 * time.Hour)

    // Config consumer
	subscription, err := initializers.JS.PullSubscribe(channel, "", nats.PullMaxWaiting(200), nats.StartTime(startTime))
	if err != nil {
		log.Fatalf("Error while subscribing to channel: %v", err)
	}

	log.Printf("Fetching messages from the last hour in channel '%s'...\n", channel)

	// Recoger los mensajes del canal (esto no los consume, sólo los muestra)
	for msg := range subscription.Chan() {
		// Mostrar el mensaje
		log.Printf("[%s] %s: %s\n", msg.Subject, msg.Header.Get("user"), string(msg.Data))
	}

	log.Println("Finished fetching messages.")
    
}