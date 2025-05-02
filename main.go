package main

import (
	"SkipperTunnelProxy/HttpServer"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := HttpServer.NewServer(8080)

	s.Router.GET("/", HttpServer.HomeHandler)

	s.Router.GET("/time", HttpServer.TimeHandler)
	s.Router.POST("/coso", HttpServer.ParsePost)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Println("Server is running on http://localhost:8080")

	//Interrupt signal
	<-stop

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")

}
