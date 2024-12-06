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
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    subscription, err := initializers.JS.Subscribe(channel, func(msg jetstream.Msg) {
        log.Printf(string(msg.Data))
    }, jetstream.Durable("durable-subscriber"), jetstream.ManualAck())

    if err != nil {
        log.Fatalf("Error al suscribirse al canal %s: %v", channel, err)
    }

    // Mantener activa la suscripción
    <-ctx.Done()
    subscription.Unsubscribe()
}

func PublishMessage(channel, user, message string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf(channel)
    fullMessage := fmt.Sprintf("[%s]: %s", user, message)

    _, err := initializers.JS.Publish(ctx, channel, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error al publicar mensaje en el canal %s: %v", channel, err)
    }
}

func FetchRecentMessages(channel string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf(channel) //("chat.%s", channel)

    startTime := time.Now().Add(-1 * time.Hour)

    // Config consumer
    consumer, err := initializers.ChatStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
        Name:          "consumer",
        Durable:       "recent-messages",
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

    msgs, err := consumer.Fetch(10)
    if err != nil {
        // handle error
    } else {
        for msg := range msgs.Messages() {
            fmt.Printf("Received a JetStream message: %s\n", string(msg.Data()))
        }
    }
    log.Println("History mss finished!")
}