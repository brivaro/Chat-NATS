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
    subject := fmt.Sprintf("chat.%s", channel)
    _, err := initializers.Client.Conn.Subscribe(subject, func(msg *nats.Msg) {
        log.Printf("Received mss: %s", string(msg.Data))
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
        Durable:       "chat_consumer",  
        FilterSubject: subject,          // Filtro para el canal específico
        DeliverPolicy: jetstream.DeliverByStartTimePolicy, // Entrega mensajes desde un momento específico
        OptStartTime:  &startTime,       // Hora de inicio (última hora)
        AckPolicy:     jetstream.AckNonePolicy, // Política de recepción
    }

    // Creating/updating consumer
    consumer, _ := initializers.JS.CreateOrUpdateConsumer(ctx, "chat_stream", consumerConfig)

    // Recovery mss
    consumer.Consume(func(msg jetstream.Msg){
        meta, err := msg.Metadata()
        if err != nil{
            log.Println("Error getting metadata:", err)
        }
        log.Println("Processing mss:", meta.Sequence.Consumer, msg.Subject())
    })
    
}