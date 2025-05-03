package main

import (
	"fmt"
	"net"
)

// thats to receive data onn a buffer and readed
func HandleClient(conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Print(err)
			return
		}

		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

func HandleSend(conn net.Conn) {

	data := []byte("Hello, Server!\n")
	_, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("sendet")
}
