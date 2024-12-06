package main

import (
	"os"
	"log"
	"strings"
	"bufio"

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

		log.Print("Channel (Format: chat>):")
		channel, _ = reader.ReadString('\n')
		channel = strings.TrimSpace(channel)

		log.Print("User: ")
		user, _ = reader.ReadString('\n')
		user = strings.TrimSpace(user)

		if serverURL != "" && channel != "" && user != "" {
			break
		} else {
			log.Println("\nAll fields are required. Please try again.\n")
		}
	}

    initializers.ConnectToNats(serverURL)
    initializers.CreateJetStream()
    initializers.CreateChatStream()

	// Subscribe
    api.SubscribeToChannel(channel)

    // Recovery historic mss
    api.FetchRecentMessages(channel)

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
