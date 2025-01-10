package utils

import (
	"fmt"
	"log"
	"net"
	"time"
)

// BroadcastMessage sends a client's message to all other clients connected to the network
// (with the exception of the sender)
func BroadcastMessage(message string, senderConnection net.Conn) {
	mClients.Lock()
	defer mClients.Unlock()

	for _, client := range clients {
		if client.connection != senderConnection { // exclude sender
			_, err := client.connection.Write([]byte(message))
			client.connection.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"),
				client.name)))
			if err != nil {
				log.Printf("Error broadcasting message to %v: %v\n", &client.name, err)
				go RemoveClient(client.connection)
			}
		}
	}
}
