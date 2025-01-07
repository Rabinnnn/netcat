package utils
import(
	"fmt"
	"time"
	"net"
	"log"
	"strings"
)

func AddNewClient(connection net.Conn){
	connection.SetDeadline(time.Now().Add(60 * time.Second))
	
	connection.Write([]byte("Welcome to TCP-Chat!\n"))
	DisplayLogo(connection)
	connection.Write([]byte("[ENTER YOUR NAME]: "))

	buffer := make([]byte, 1024)
	num, err := connection.Read(buffer)
	if err != nil{
		log.Printf("Error reading client name: %v\n", err)
		connection.Close()
		return
	}

	connection.SetDeadline(time.Time{})

	clientName := strings.TrimSpace(string(buffer[:num]))
	//fmt.Println(clientName)
	if clientName == ""{
		connection.Write([]byte("Name must not be empty\n"))
		connection.Close()
		return
	}

	mClients.Lock()
	for _, client := range clients{
		if strings.EqualFold(client.name, clientName){
			mClients.Unlock()
			connection.Write([]byte("The name has already been taken by another user\n"))
			connection.Close()
			return
		}
	}
	mClients.Unlock()

	newClient := &client{
		name: clientName,
		connection: connection,
	}
	mClients.Lock()
	clients = append(clients, newClient)
	mClients.Unlock()

	DisplayChats(newClient)
	BroadcastMessage(fmt.Sprintf("\n%v has joined our chat...", clientName), connection)
	go HandleClientSession(newClient)

}