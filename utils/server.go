package utils

import(
	"log"
	"net"
	"sync"
	"strings"
	"strconv"
)

type client struct {
	name string
	connection net.Conn
}

var maxConnections = 10 // maximum number of connected clients allowed.
var clients = make([]*client, 0, maxConnections) // slice holding all connected clients
var mClients sync.Mutex // mutex to synchronize access to clients slice
var chatHistory []string // slice to store messages
var mChatHistory sync.Mutex // mutex to synchronize access to chatHistory

// Server sets up a tcp server for the application.
// It:
// 	- listens on the specified port and accept incomming client connections
//	- gracefully handles errors that might arise during the process
// 	- ensures the number of connected clients doesn't exceed 10

func Server(port string){
	portNum, err := strconv.Atoi(strings.TrimPrefix(port, ":"))
	if err != nil{
		log.Printf("Invalid port number %q: %v\n", port, err)
		return
	}

	listener, err := net.Listen("tcp", port)
	if err != nil{
		log.Printf("Error listening on the port %d: %v\n", portNum, err)
		return
	}
	defer listener.Close()
	
	log.Printf("Listening on the port :%d\n", portNum)
	//fmt.Printf("Listening on port %q\n", port)

	for {
		connection, err := listener.Accept()
		if err != nil{
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		mClients.Lock()
		if len(clients) >= maxConnections {
			mClients.Unlock()
			connection.Write([]byte("The group chat is currently full. Please try again later.\n"))
			connection.Close()
			continue
		}
		mClients.Unlock()

		go AddNewClient(connection)
	}

}