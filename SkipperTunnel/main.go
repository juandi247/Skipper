package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	requestChannel := make(chan string)
	go func() {
		NewHttpClient()
		fmt.Println("inicando HTTP")

	}()

	
	go func() {
		conn, err := net.Dial("tcp", "localhost:9000")
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Println("estoy inciando TCP")

		// todo: make them a goruoitne and wrap the entire tcp connection on a gorountie

		go HandleReceive(conn, requestChannel)
		go ReceiveRequest(requestChannel)
		defer conn.Close()

		for {
			time.Sleep(10 * time.Second)
			fmt.Println("Conexión TCP sigue activa")
		}

	}()
	select {}

}
