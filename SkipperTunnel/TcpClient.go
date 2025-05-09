package main

import (
	"fmt"
	"io"
	"net"
)

// thats to receive data onn a buffer and readed
func HandleReceive(conn net.Conn, ch chan string) {
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
			ch <- string(buffer[:n])
		}
	}
}

func HandleSend(conn net.Conn) {
	data := []byte("Hello, Server!\n")
	_, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("sendet Hello server message")
}
