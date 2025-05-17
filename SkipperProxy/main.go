package main

import (
	"SkipperTunnelProxy/HttpServer"
	tcpserver "SkipperTunnelProxy/TcpServer"
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
	server := tcpserver.NewServer(":9000")
	// Run http server
	s := HttpServer.NewServer(8080, server.RequestChannel, false)

	// ! just for prod enviroment on GCP virtual machine
	// httpsServer:= HttpServer.NewServer(443, server.RequestChannel, true)

	s.Router.Any("/*", s.ParseHttpRequest)
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
		for msg := range server.MessageChanel {
			// Intentar deserializar el mensaje como JSON
			var httpResponse HttpServer.HttpResponse
			err := json.Unmarshal(msg, &httpResponse)
			if err != nil {
				fmt.Println("❌ Error al deserializar el mensaje:", err)
				continue
			}

			// Verificar si el JSON tiene los campos esperados
			fmt.Println("Cuerpo:", httpResponse.Body)
			fmt.Println("IDDDD:", httpResponse.RequestID)

			// Obtener el requestID de los datos si lo necesitas
			requestID := httpResponse.RequestID // Asumir que 'Method' es el requestID, o usa otro campo
			responseChannel, ok := s.GetResponseChannel(requestID)
			if ok {
				// Enviar el cuerpo del mensaje al canal de respuesta
				responseChannel <- []byte(httpResponse.Body)
			} else {
				fmt.Printf("⚠️ No se encontró canal de respuesta para el requestID: %s\n", requestID)
			}
		}
	}()

	// this is for the starting of the server
	go server.Start()

	go func() {
		for msg := range server.RequestChannel {
			server.ConnMutext.Lock()

			fmt.Printf("Buscando conexión para target: %s\n", msg.Target)
			fmt.Printf("Contenido de ConnectionMap: %v\n", server.ConnectionMap)

			conn, exists := server.ConnectionMap[msg.Target]
			fmt.Println(exists)
			if !exists {
				fmt.Printf("⚠️ No se encontró conexión para el target: %s\n", msg.Target)
				server.ConnMutext.Unlock()
				continue
			}
			_, err := conn.Write(msg.Data)
			if err != nil {
				fmt.Printf("❌ Error escribiendo a conexión TCP [%s]: %v\n", msg.Target, err)
				delete(server.ConnectionMap, msg.Target)
				conn.Close()
			} else {
				fmt.Printf("✅ Mensaje enviado a [%s]\n", msg.Target)
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
