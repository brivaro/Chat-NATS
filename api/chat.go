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
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf("chat.%s", channel)

    startTime := time.Now().Add(-1 * time.Hour)

    // Config consumer
    consumerConfig := jetstream.ConsumerConfig{
        Durable:       "recent-msgs",  
        FilterSubject: subject,          // Filtro para el canal específico
        DeliverPolicy: jetstream.DeliverByStartTimePolicy, // Entrega mensajes desde un momento específico
        OptStartTime:  &startTime,       // Hora de inicio (última hora)
        AckPolicy:     jetstream.AckExplicitPolicy, // Política de recepción
        ReplayPolicy:   jetstream.ReplayInstantPolicy, // Política de Reproducción
    }

    // Creating/updating consumer
    _, err := initializers.JS.CreateOrUpdateConsumer(ctx, "chatSAD", consumerConfig)
    if err != nil {
        log.Fatalf("Error creating consumer: %v", err)
    }

    sub, err := initializers.Client.Conn.PullSubscribe(subject, "recent-msgs")
    if err != nil {
        log.Fatalf("Error subscribing: %v", err)
    }

    go func() {
        for {
            msgs, err := sub.Fetch(10, nats.Context(context.Background()))
            if err != nil {
                log.Printf("Error fetching messages: %v", err)
                continue
            }
            for _, msg := range msgs {
                log.Printf(string(msg.Data))
                msg.Ack()
            }
        }
    }()
    
}