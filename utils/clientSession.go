package utils

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// HandleClientSession manages interaction between a connected client and the chat server.
func HandleClientSession(client *client) {
	defer RemoveClient(client.connection)
	buffer := make([]byte, 4096)

	for {
		client.connection.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), client.name)))
		num, err := client.connection.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				BroadcastMessage(fmt.Sprintf("\n%s has left our chat...", client.name), client.connection)
				return
			}
			log.Printf("Error reading from %v: %v\n", client.name, err)
			return
		}

		clientMessage := strings.TrimSpace(string(buffer[:num]))
		if clientMessage == "" {
			log.Println("Can't send an empty message to the chat.")
			continue
		}

		formattedClientMessage := ""
		mChatHistory.Lock()
		formattedClientMessage = fmt.Sprintf("\n[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), client.name, clientMessage)
		
		if len(chatHistory) == 0{
			chatHistory = append(chatHistory, fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), client.name, clientMessage))
		}else{
			chatHistory = append(chatHistory, formattedClientMessage)
		}
		mChatHistory.Unlock()

		BroadcastMessage(formattedClientMessage, client.connection)
	}
}
