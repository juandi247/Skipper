package TcpUserClient

import (
	"fmt"
	"io"
	"net"
)

// thats to receive data onn a buffer and readed
func HandleReceive(conn net.Conn, ch chan []byte) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("El servidor cerró la conexión.")
				return
			}
			fmt.Print(err)
			return
		}
		fmt.Printf("Received: %s\n", buffer[:n])
		// Si se reciben datos, envíalos al canal
		if n > 0 {
			fmt.Printf("Received: %s\n", buffer[:n])
			ch <- buffer[:n]
		}
	}
}

func HandleSend(ch chan []byte, conn net.Conn) {
	for {
		response := <-ch
		fmt.Println("VOY A ENVIAR", response)

		data := []byte(response)
		_, err := conn.Write(data)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("sendet Hello server message")
	}
}
