/*
Copyright © 2025 Juan Diego Diaz <juand.diaza@gmail.com>
*/
package cmd

import (
	"SkipperTunnel/HttpUserClient"
	"SkipperTunnel/TcpUserClient"
	"SkipperTunnel/utils"

	"fmt"
	"net"
	"os"
	"runtime"
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

		clientHttp := HttpUserClient.NewHttpCliennt(localhostUrl, 5*time.Second)
		// HTTP CLIENTT
		go func() {
			resp, err := clientHttp.Client.Get(localhostUrl)
			if err != nil {
				fmt.Println("Error on localhost", err)
				os.Exit(1)
				return
			}
			// go HttpUserClient.ReceiveRequest(requestChannel, clientHttp)
			// todo: worker pool to limit the goroutines inside, and make kind of a queue of gouroitnes

			go HttpUserClient.ReceiveRequest(requestChannel, responseChannel, clientHttp)

			defer resp.Body.Close()

			for {
				respPing, err := utils.Ping(localhostUrl, clientHttp.Client)
				if err != nil || respPing != 200 {
					fmt.Println("ping completed to localhost")
					return
				}
				time.Sleep(10 * time.Second)
				fmt.Println("ping completed to localhost")
			}
		}()

		// TCP CLIENT
		go func() {
			conn, err := net.Dial("tcp", proxyUrl)
			if err != nil {
				fmt.Print(err)
				return
			}
			fmt.Println("estoy inciando TCP")
			i, err := conn.Write([]byte(subdomain))

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("sendet Hello server message", i)

			go TcpUserClient.HandleReceive(conn, requestChannel)

			go TcpUserClient.HandleSend(responseChannel, conn)

			// go TcpUserClient.HandleReceive(conn, requestChannel)
			defer conn.Close()

			for {
				time.Sleep(2 * time.Second)
				runtime.GOMAXPROCS(0) // Usar todos los núcleos del CPU
				fmt.Println("Número total de goroutines:", runtime.NumGoroutine())
			}
		}()

		select {}

	},
}

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
