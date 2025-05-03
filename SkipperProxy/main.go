package main

import (
	"SkipperTunnelProxy/HttpServer"
	tcpserver "SkipperTunnelProxy/TcpServer"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := tcpserver.NewServer(":9000")
	// Run http server
	s := HttpServer.NewServer(8080, server.RequestChannel)
	s.Router.Any("/*", s.ParseHttpRequest)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// ! Goroutine for running the server
	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Println("Server is running on http://localhost:8080")

	// ! Run TCP server
	go func() {
		for msg := range server.MessageChanel {
			fmt.Println("message received", string(msg))
		}
	}()
	go server.Start()

	go func() {
		for msg := range server.RequestChannel {
			server.ConnMutext.Lock()
			for _, conn := range server.ConnectionMap {
				// Solo uno en el map, le escribimos y salimos
				_, err := conn.Write(msg.Data)
				if err != nil {
					fmt.Printf("Error escribiendo a conexión TCP: %v\n", err)
				}
				break 
			}
			server.ConnMutext.Unlock()
		}
	}()

	// STOPPPP the http when getting the STOP
	//Interrupt signal
	<-stop
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
