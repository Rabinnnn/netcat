package utils

import (
	"net"
)

func DisplayLogo(connection net.Conn, filepath string) {
	connection.Write([]byte(GetLogo(filepath)))
}
