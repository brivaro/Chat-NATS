package main

import (
	"os"
	"log"
	"strings"
	"bufio"
	"os/signal"
	"time" 
    "context" 

	"chat/api"
	"chat/initializers"
)

func main() {
	var serverURL, channel, user string

	reader := bufio.NewReader(os.Stdin)

	for {
		log.Println("Please, insert the next values to join a chat:")
		log.Print("Server URL: ")
		serverURL, _ = reader.ReadString('\n')
		serverURL = strings.TrimSpace(serverURL)

		log.Print("Channel (Format: chat.>):") // nats://localhost:4222
		channel, _ = reader.ReadString('\n')
		channel = strings.TrimSpace(channel)

		log.Print("User (Without: whitespace, ., *, >, path separators (forward or backward slash), or non-printable characters): ")
		user, _ = reader.ReadString('\n')
		user = strings.TrimSpace(user)

		if serverURL != "" && channel != "" && user != "" {
			break
		} else {
			log.Println("\nAll fields are required. Please try again.\n")
		}
	}

    initializers.ConnectToNats(serverURL)
	defer initializers.Client.Conn.Drain()

    initializers.CreateJetStream()
    initializers.CreateChatStream()

	// Subscribe
    api.SubscribeToChannel(channel, user)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Goroutine CTRL+C to clean resources
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		log.Println("Shutting down... cleaning up resources.")
		// Eliminar consumidor antes de salir
		if api.Consumer != nil {
			if err := initializers.ChatStream.DeleteConsumer(ctx, api.Consumer.CachedInfo().Name); err != nil {
				log.Printf("Error deleting consumer: %v", err)
			} else {
				log.Println("Consumer deleted successfully. You exit the channel!")
			}
		}
		cancel()
		os.Exit(0)
	}()

    // Read terminal mss
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        if strings.TrimSpace(text) != "" {
            api.PublishMessage(channel, user, text)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading input: %v", err)
    }
}

