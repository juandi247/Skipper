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
	"sync"

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
// const proxyUrl string = "localhost:9000"

// !prod
const proxyUrl string = "skipper.lat:80"

// startskipperCmd represents the startskipper command
var startskipperCmd = &cobra.Command{
	Use:   "skipper start",
	Short: "This command starts the tunnel connection and lets you expose your localhost on the web [your-subdomain].skipper.lat,\n to use this command please use the -p (port) flag and the -s (subdomain) flag",
	Run: func(cmd *cobra.Command, args []string) {
		requestChannel := make(chan []byte, 30)
		localhostUrl := "http://localhost:" + strconv.Itoa(port)
		ctx, gracefullShutdown := context.WithCancel(context.Background())

		clientHttp := HttpUserClient.NewHttpCliennt(localhostUrl, 5*time.Second)
		// HTTP CLIENTT connection
		for i := 0; i < 5; i++ {
			fmt.Println("Trying to stablish connection with localhost:", port)
			resp, err := clientHttp.Client.Get(localhostUrl)
			if err == nil {
				resp.Body.Close()
				break
			} else if err != nil && i == 5 {
				fmt.Printf("Could not find an active localhost on port %d, %v ", port, err)
				gracefullShutdown()
				// os.Exit(1)
				return
			}
			time.Sleep(2 * time.Second)
		}

		fmt.Println("Connection Stablished Succesfully")

		// TCP CONNECTION HANDLER
		conn, err := net.Dial("tcp", proxyUrl)
		if err != nil {
			fmt.Print("Error connecting to Skippers Proxy", err)
			gracefullShutdown()
			return
		}
		defer conn.Close()

		i, err := conn.Write([]byte(subdomain))

		if err != nil {
			fmt.Println("An error ocurred", err, i)
			gracefullShutdown()
			return
		}
		fmt.Printf("You can now use the https://%s.skipper.lat subdomain \n", subdomain)

		// ! GOROUTINES
		/* wait group (this was made because of the handling
		receive from the TCP proxy, the io.ReadFull is a blocking method
		so using the waitgroup allows us to make use of the context better and then
		make a good gracefull shutdown without using an empty select{} on the main function) */
		var wg sync.WaitGroup
		// ping localhost goroutine
		wg.Add(17)
		go func(ctx context.Context, wg *sync.WaitGroup) {
			for {
				respPing, err := utils.Ping(localhostUrl, clientHttp.Client)
				if err != nil || respPing != 200 {
					// fmt.Println("ping ME FALLO to localhost")
					gracefullShutdown()
					wg.Done()
					// os.Exit(1)
					return
				}
				// fmt.Println("ping completed to localhost")
				select {
				case <-time.After(3 * time.Second): // Espera 3 segundos
				case <-ctx.Done(): // Si el contexto se cancela durante el sleep, sal inmediatamente
					// fmt.Println("Goroutine de ping: Contexto cancelado durante el sleep. Terminando.")
					wg.Done()
					return
				}
			}
		}(ctx, &wg)
		// goroutine for handling tcp connection
		go TcpUserClient.HandleReceive(conn, requestChannel, ctx, &wg)

		// ! worker pool to handle the requests
		for i := 0; i < 15; i++ {
			go HttpUserClient.ReceiveRequest(localhostUrl, i, requestChannel, clientHttp, conn, ctx, &wg)
		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(startskipperCmd)

	startskipperCmd.Flags().IntVarP(&port, "port", "p", 8080, "Your Localhost Port")
	startskipperCmd.Flags().StringVarP(&subdomain, "subdomain", "s", "", "Subdomain that you want")

	if err := startskipperCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println("Error: Flag -p or port not founded: ", err)
		os.Exit(1)
	}
	if err := startskipperCmd.MarkFlagRequired("subdomain"); err != nil {
		fmt.Println("Error: flag -s or subdomain not founded", err)
		os.Exit(1)
	}
}
