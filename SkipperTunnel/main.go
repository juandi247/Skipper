package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("inicando conex")
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("bien hasta ahora")
	handleSend(conn)
	time.Sleep(3 * time.Second)
	handleClient(conn)
	defer conn.Close()

}

// thats to receive data onn a buffer and readed
func handleClient(conn net.Conn) {
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

func handleSend(conn net.Conn) {

	data := []byte("Hello, Server!\n")
	_, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("sendet")
}
