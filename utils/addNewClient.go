package utils

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func AddNewClient(connection net.Conn) {
	connection.SetDeadline(time.Now().Add(60 * time.Second))

	connection.Write([]byte("Welcome to TCP-Chat!\n"))
	DisplayLogo(connection)

	for {
		connection.Write([]byte("[ENTER YOUR NAME]: "))
		buffer := make([]byte, 1024)
		num, err := connection.Read(buffer)
		if err != nil {
			log.Printf("Error reading client name: %v\n", err)
			connection.Close()
			return
		}

		clientName := strings.TrimSpace(string(buffer[:num]))

		if clientName == "" {
			connection.Write([]byte("Name must not be empty. Please try again.\n"))
			continue
		}

		mClients.Lock()
		isDuplicate := false
		for _, client := range clients {
			if strings.EqualFold(client.name, clientName) {
				isDuplicate = true
				break
			}
		}
		mClients.Unlock()

		if isDuplicate {
			connection.Write([]byte("The name has already been taken by another user. Please choose a different name.\n"))
			continue
		}

		connection.SetDeadline(time.Time{}) // Remove deadline after successful name entry

		newClient := &client{
			name:       clientName,
			connection: connection,
		}

		mClients.Lock()
		clients = append(clients, newClient)
		mClients.Unlock()

		DisplayChats(newClient)
		BroadcastMessage(fmt.Sprintf("\n%v has joined our chat...", clientName), connection)
		go HandleClientSession(newClient)
		return
	}
}
