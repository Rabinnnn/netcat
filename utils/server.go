package utils

import(
	"fmt"
	"log"
	"net"
	"sync"
)

type client struct {
	name string
	connection net.Conn
}

var maxConnections = 10
var clients = make([]*client, 0, maxConnections)
var mAll sync.Mutex
var chatHistory []string
var mChatHistory sync.Mutex

func Server(port string){
	listener, err := net.Listen("tcp", port)
	if err != nil{
		log.Printf("Error listening on port %q: %v\n", port, err)
		return
	}
	defer listener.Close()

	log.Printf("Listening on port %q\n", port)
	fmt.Printf("Listening on port %q\n", port)

	for {
		connection, err := listener.Accept()
		if err != nil{
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		mAll.Lock()
		if len(clients) >= maxConnections {
			mAll.Unlock()
			connection.Write([]byte("The group chat is currently full. Please try again later.\n"))
			connection.Close()
			continue
		}
		mAll.Unlock()
	}

}