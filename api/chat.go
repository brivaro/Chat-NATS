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
    fullMessage := fmt.Sprintf("[%s]: %s", user, message)
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
    consumer, err := initializers.ChatStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
        Name:          "recent-messages",
        Durable:       "",
        Description:   "Consumer to fetch recent messages",
        FilterSubject: subject,                             // Filtro para el canal específico
        DeliverPolicy: jetstream.DeliverByStartTimePolicy, // Inicia desde una hora específica
        OptStartTime:  &startTime,                        // Hora de inicio
        AckPolicy:     jetstream.AckNonePolicy,           // Política de recepción: No consume los mensajes
        ReplayPolicy:  jetstream.ReplayInstantPolicy,     // Política de Reproducción: Reproducción instantánea
    })
    // Creating/updating consumer
    if err != nil {
        log.Fatalf("Error creando consumidor: %v", err)
    }
    log.Println("Consumer created to recovery history mss.")

    _, err = consumer.Consume(func(msg jetstream.Msg) {
        fmt.Printf("[%s]: %s\n", 
            msg.Header.Get("user"), 
            string(msg.Data),
        )

        // No se envía ACK para mantener los mensajes disponibles
    })
    if err != nil {
        log.Fatalf("Error consuming mss: %v", err)
    }

    log.Println("History mss finished!")
}