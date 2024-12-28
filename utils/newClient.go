package utils
import(
	"time"
	"net"
)

func clientConnection(connection net.Conn){
	connection.SetDeadline(time.Now().Add(60 * time.Second))
	
	DisplayLogo(connection)
	connection.Write([]byte("[ENTER YOUR NAME]: "))
}