package utils

import (
	"net"
)

func DisplayLogo(connection net.Conn) {
	connection.Write([]byte(GetLogo("linuxLogo.txt")))
}
