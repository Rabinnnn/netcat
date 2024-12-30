package utils
import(
	"fmt"
	"net"
	"log"
	"time"	
)

// BroadcastMessage sends a client's message to all other clients connected to the network
// (with the exception of the sender)
func BroadcastMessage(message string, senderConnection net.Conn){
	mChatHistory.Lock()
	chatHistory = append(chatHistory, message)
	mChatHistory.Unlock()

	mClients.Lock()
	defer mClients.Unlock()

	for _, client := range clients{
		if client.connection != senderConnection{ // exclude sender
			_, err := client.connection.Write([]byte(message))
			client.connection.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(time.DateTime),
			client.name)))
			if err != nil{
				log.Printf("Error broadcasting message to %v: %v\n", &client.name, err)
				go RemoveClient(client.connection)
			}
		}
	}
}