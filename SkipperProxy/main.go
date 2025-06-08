package main

import (
	"SkipperProxy/HttpServer"
	tcpserver "SkipperProxy/TcpServer"
	"SkipperProxy/config"
	"SkipperProxy/connectionmanager"
	"SkipperProxy/worker"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)
// var for enviroment setting
var Env string

func main() {
	fmt.Println("ENV:", Env)

	config := config.LoadConfig(Env)
	fmt.Println("ENV:", Env)

	fmt.Println("My base domain is", config.BaseDomain)

	cm := connectionmanager.NewConnectionManager()
	tcpserver := tcpserver.NewServer(config.TcpPort, cm)
	// Run http server
	s := HttpServer.NewServer(config.HttpPort, config.IsHttps, cm, config)
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error cargando templates: %v", err)
		return
	}
	s.Templates = templates

	s.Router.Any("/*", s.HandleClientRequest)
	s.Router.ServeFavicon()
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

	wpChannel := make(chan []byte, 50)
	go func() {
		for msg := range tcpserver.MessageChanel {
			wpChannel <- msg
		}
	}()

	// worker pool
	// todo check worker pool size because of low specificacions of ram on the VM
	for i := 0; i < config.WorkerNumber; i++ {
		fmt.Println("creating gorotounie", i)
		go worker.StartWorker(i, wpChannel, cm)
	}

	go tcpserver.Start()

	// go func() {
	// 	for msg := range tcpserver.RequestChannel {
	// 		tcpserver.ConnMutext.Lock()
	// 		for _, conn := range tcpserver.ConnectionMap {
	// 			// Solo uno en el map, le escribimos y salimos
	// 			_, err := conn.Write(msg.Data)
	// 			fmt.Println("ya se fue")
	// 			if err != nil {
	// 				fmt.Printf("Error escribiendo a conexión TCP: %v\n", err)
	// 			}
	// 			break
	// 		}
	// 		tcpserver.ConnMutext.Unlock()
	// 	}
	// }()

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
