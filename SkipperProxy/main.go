package main

import (
	"SkipperProxy/http"
	"SkipperProxy/tcp"
	"SkipperProxy/tunnel"
	"fmt"
	"sync"
)

/*
this is the new skipper rewrite implementation, using new features and implementing
TigerStyle inspiration for the project. Im not expert or anything in go and programming but this
project is focused for learning Golang mainly and experimenting with new low level concepts :)
Juan Diego Diaz
*/
func main() {
	wg := sync.WaitGroup{}
	httpServer := http.CreateHttpServer()
	wg.Add(1)
	go httpServer.StartServer()
	tm := tunnel.CreateTunnelManager()
	tcpServer := tcp.CreateTcpServer(":9000", tm)
	wg.Add(1)
	go tcpServer.StartServer()
	fmt.Println("tcp server started")
	wg.Wait()

}
