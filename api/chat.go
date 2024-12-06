package api

import (
	"fmt"        
	"log"        
	"time" 
    "context"  
	"chat/initializers" 
    "strconv"
    "math/rand"

	"github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

func SubscribeToChannel(channel string) {
    _, err := initializers.Client.Conn.Subscribe(channel, func(msg *nats.Msg) {
        log.Printf(string(msg.Data))
    })

    if err != nil {
        log.Fatalf("Error al suscribirse al canal %s: %v", channel, err)
    }
}

func PublishMessage(channel, user, message string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf(channel)
    fullMessage := fmt.Sprintf("[%s]: %s", user, message)

    err := initializers.Client.Conn.Publish(subject, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error al publicar mensaje en el canal %s: %v", channel, err)
    }

    _, err = initializers.JS.Publish(ctx, subject, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error al publicar mensaje en el canal %s: %v", channel, err)
    }
}

func FetchRecentMessages(channel string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    subject := fmt.Sprintf(channel) //("chat.%s", channel)

    //startTime := time.Now().Add(-1 * time.Hour)
    randId := rand.Intn(10) + 1

    // Config consumer
    consumer, err := initializers.ChatStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
        Name:          fmt.Sprintf("consumer_%s", strconv.Itoa(randId)),
        Durable:       "recent-messages",
        Description:   "Consumer to fetch recent messages",
        FilterSubject:  subject,                             // Filtro para el canal específico
        //DeliverPolicy: jetstream.DeliverByStartTimePolicy, // Inicia desde una hora específica
        //OptStartTime:  &startTime,                        // Hora de inicio
        //AckPolicy:     jetstream.AckNonePolicy,           // Política de recepción: No consume los mensajes
        //ReplayPolicy:  jetstream.ReplayInstantPolicy,     // Política de Reproducción: Reproducción instantánea
    })

	log.Printf("Created consumer")
	if err != nil {
		log.Fatal(err)
	}
    
    log.Println("Created consumer", consumer.CachedInfo().Name)
}