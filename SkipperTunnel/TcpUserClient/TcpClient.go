package TcpUserClient

import (
	"bytes"
	"encoding/binary"
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
		lenght := uint32(len(response))

		buf := new(bytes.Buffer)

		binary.Write(buf, binary.BigEndian, lenght)
		buf.Write(response)

		// todo: check becasue we are reconverting the resopne to bytes, and the resopnse was already in bytes
		// data := []byte(response)
		_, err := conn.Write(buf.Bytes())

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("sendet Hello server message")
	}
}
