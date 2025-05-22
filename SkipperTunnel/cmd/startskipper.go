/*
Copyright © 2025 Juan Diego Diaz <juand.diaza@gmail.com>
*/
package cmd

import (
	"SkipperTunnel/HttpUserClient"
	"SkipperTunnel/TcpUserClient"
	"SkipperTunnel/utils"
	"context"
	"fmt"
	"net"
	"os"

	// "runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	port      int
	subdomain string
)

// ? dev
const proxyUrl string = "localhost:9000"

// !prod
// const proxyUrl string = "skipper.lat:80"

// startskipperCmd represents the startskipper command
var startskipperCmd = &cobra.Command{
	Use:   "startskipper",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		requestChannel := make(chan []byte)
		responseChannel := make(chan []byte)
		localhostUrl := "http://localhost:" + strconv.Itoa(port)
		ctx, gracefullShutdown := context.WithCancel(context.Background())

		clientHttp := HttpUserClient.NewHttpCliennt(localhostUrl, 5*time.Second)
		// HTTP CLIENTT connection
		for i := 0; i < 5; i++ {
			fmt.Println("trying connection")
			resp, err := clientHttp.Client.Get(localhostUrl)
			if err==nil{
				break
			}else if err!=nil && i==5{
				fmt.Printf("Could not find an active localhost on port %d, %v ", port, err)
				gracefullShutdown()
				resp.Body.Close()
				// os.Exit(1)
				return
			} 
			resp.Body.Close()
			time.Sleep(2*time.Second)
		}
	
		// TCP CONNECTION HANDLER
		conn, err := net.Dial("tcp", proxyUrl)
		if err != nil {
			fmt.Print("ERRORsote", err)
			gracefullShutdown()
			return
		}
		defer conn.Close()

		fmt.Println("estoy inciando TCP")
		i, err := conn.Write([]byte(subdomain))

		if err != nil {
			fmt.Println("Error:", err)
			gracefullShutdown()
			return
		}
		fmt.Println("sendet Hello server message", i)


	// TODO: handle correctly with WaitGroup for a better cancelation and with a general select 
	// TODO: error handling channel for child goroutines
		// ! GOROUTINES
		// ping localhost goroutine
		go func(ctx context.Context) {
			for {
				respPing, err := utils.Ping(localhostUrl, clientHttp.Client)
				if err != nil || respPing != 200 {
					fmt.Println("ping ME FALLO to localhost")
					gracefullShutdown()
					// os.Exit(1)
					return
				}
				fmt.Println("ping completed to localhost")
				select {
				case <-time.After(3 * time.Second): // Espera 3 segundos
				case <-ctx.Done(): // Si el contexto se cancela durante el sleep, sal inmediatamente
					fmt.Println("Goroutine de ping: Contexto cancelado durante el sleep. Terminando.")
					return
				}
			}
		}(ctx)
		// goroutine for handling tcp connection
		go TcpUserClient.HandleReceive(conn, requestChannel, ctx)
		// goruotine to handle send httprequest to localhost and send response using TCP
		go HttpUserClient.ReceiveRequest(requestChannel, responseChannel, clientHttp, conn, ctx)

		select {}

	},
}

// for {
// 	time.Sleep(2 * time.Second)
// 	runtime.GOMAXPROCS(0) // Usar todos los núcleos del CPU
// 	fmt.Println("Número total de goroutines:", runtime.NumGoroutine())
// }
// }()

func init() {
	rootCmd.AddCommand(startskipperCmd)

	startskipperCmd.Flags().IntVarP(&port, "port", "p", 8080, "Your Localhost Port")
	startskipperCmd.Flags().StringVarP(&subdomain, "subdomain", "s", "", "Subdomain that you want")

	if err := startskipperCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println("Error: el flag -p es obligatorio (Puerto)", err)
		os.Exit(1)
	}
	if err := startskipperCmd.MarkFlagRequired("subdomain"); err != nil {
		fmt.Println("Error: el flag -s es obligatorio (Subdominio)", err)
		os.Exit(1)
	}
}
