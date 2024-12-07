package api

import (
	"fmt"        
	"log"        
	"time" 
    "context"  
	"chat/initializers" 
    //"strconv"
    //"math/rand"

	//"github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

var ConsumerCon jetstream.ConsumeContext
var Consumer jetstream.Consumer

func SubscribeToChannel(channel string, user string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    startTime := time.Now().Add(-1 * time.Hour)
    //randId := rand.Intn(1000) + 1
    c := fmt.Sprintf("consumer_%s", user)

    // Check if a consumer with the same name already exists
    _, err := initializers.ChatStream.Consumer(ctx, c)
    if err == nil { // No error means the consumer exists
        log.Fatalf("A consumer with the same name already exists.")
    } else if err != jetstream.ErrConsumerNotFound {
        log.Fatalf("Error checking consumer: %v", err)
    }
    
    consumer, err := initializers.ChatStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
        Name:          c,
        Durable:       c,
        Description:   c,
        FilterSubject:  channel,                             // Filtro para el canal específico
        //InactiveThreshold: 10 * time.Millisecond,
        DeliverPolicy: jetstream.DeliverByStartTimePolicy, // Inicia desde una hora específica
        OptStartTime:  &startTime,                        // Hora de inicio
        AckPolicy:     jetstream.AckExplicitPolicy,           // Política de recepción: No consume los mensajes
        ReplayPolicy:  jetstream.ReplayInstantPolicy,     // Política de Reproducción: Reproducción instantánea
	})
    if err != nil {
        log.Fatalf("Error creating consumer: %v", err)
    }
	fmt.Println("Created consumer", consumer.CachedInfo().Name)

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		fmt.Println(string(msg.Data()))
		msg.Ack()
	})
    if err != nil {
        log.Fatalf("Error creating consumerContext: %v", err)
    }
    ConsumerCon = cc
    Consumer = consumer
}

func PublishMessage(channel, user, message string) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    fullMessage := fmt.Sprintf("%s [%s]: %s", time.Now().Format("15:04:05"), user, message)

    _, err := initializers.JS.Publish(ctx, channel, []byte(fullMessage))
    if err != nil {
        log.Fatalf("Error publishing mss to the channel %s: %v", channel, err)
    }
}

