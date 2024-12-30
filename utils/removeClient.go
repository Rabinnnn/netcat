package utils
import(
	"net"
)

// RemoveClient removes a specific client from the list of connected clients
// It gets the correct client by checking the client's connection that matches
// the passed connection.
func RemoveClient(connection net.Conn){
	mClients.Lock()
	defer mClients.Unlock()

	for i, client := range clients{
		if client.connection == connection{
			client.connection.Close()
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}