package main

import (
	"log"
	"os"
	"strconv"

	"netcat/utils"
)

const defaultPort = ":8989"

func main() {
	switch len(os.Args) {
	case 1:
		utils.Server(defaultPort)
	case 2:
		port, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Printf("Error converting %q to an int: %v\n", port, err)
			return
		}

		if port < 1024 || port > 65535 {
			log.Println("Invalid port. Allowed range is 1024 - 65535")
			return
		}
		utils.Server(":" + strconv.Itoa(port))
	default:
		log.Println("[USAGE]: ./TCPChat $port")

	}
}
