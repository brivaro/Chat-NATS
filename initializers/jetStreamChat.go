package initializers

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

var ChatStream jetstream.Stream

var streamConfig = jetstream.StreamConfig{
	Name:       "SAD",                 // Nombre del stream
	Subjects:   []string{"chat>"},         // Sujetos que el stream observará
	Storage:    jetstream.FileStorage,     // Tipo de almacenamiento: archivo
	Replicas:   1,                         // Número de réplicas
	Retention:  jetstream.LimitsPolicy,    // Retiene mensajes según los límites // Política de retención
	Discard:    jetstream.DiscardOld,      // Descartar mensajes más antiguos si se supera el límite
	//MaxMsgs:    1000,                     // Límite de mensajes en el stream
	MaxBytes:   128 * 1024 * 1024,         // Tamaño máximo total del stream (128MB)
	MaxAge:     2 * time.Hour,             // Tiempo de vida (TTL) de los mensajes (1 hour)
	MaxMsgSize: -1,                        // Sin límite en el tamaño de mensajes
}

func CreateChatStream() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := JS.CreateStream(ctx, streamConfig)
	if err != nil {
		if err == jetstream.ErrStreamNameAlreadyInUse {
			//log.Printf("Chat already exists!")
		} else {
			log.Fatal(err)
		}
	}
	ChatStream = stream
	log.Printf("Welcome to the chat! Type your message below ;)\n")
}