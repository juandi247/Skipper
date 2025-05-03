package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	NewHttpClient()

	fmt.Println("inicando conex")
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("bien hasta ahora")
	HandleSend(conn)
	time.Sleep(3 * time.Second)
	HandleClient(conn)
	defer conn.Close()

}
