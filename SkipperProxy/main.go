package main

import (
	"SkipperTunnelProxy/HttpServer"
	tcpserver "SkipperTunnelProxy/TcpServer"
	"SkipperTunnelProxy/connectionmanager"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cm := connectionmanager.NewConnectionManager()
	tcpserver := tcpserver.NewServer(":9000", cm)

	// Run http server
	s := HttpServer.NewServer(8080, false, cm)

	// ! just for prod enviroment on GCP virtual machine
	// httpsServer:= HttpServer.NewServer(443, server.RequestChannel, true)

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
	go func() {
		for msg := range tcpserver.MessageChanel {
			fmt.Println("message received", string(msg))
			var response HttpServer.HttpResponse
			err := json.Unmarshal(msg, &response)
			if err != nil {
				fmt.Println("error parsing resopnee")
				continue
			}

			fmt.Println("le vamos a enviarrrrr")


			cm.Mu.Lock()
			ch, exists := cm.GlobalResponseChannel[response.RequestID]
			
			fmt.Println("REQUEST ID FROM RESOPNSEEEE", response.RequestID)
			if exists {
				// Enviar la respuesta al channel que espera el HTTP handler
				ch <- msg

				fmt.Println("si existio mensaje le enviamoss!!")

				// Opcional: cerrar el channel y borrarlo del mapa para limpiar
				close(ch)
				delete(cm.GlobalResponseChannel, response.RequestID)
			}
			fmt.Println("QUE PASO NO PASO NADAAAAAA")
			cm.Mu.Unlock()
		}
	}()

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
