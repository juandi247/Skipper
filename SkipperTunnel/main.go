/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"SkipperTunnel/constants"
	skipperflag "SkipperTunnel/flags"
	"SkipperTunnel/frame"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	localhostUrl string
	proxyUrl     string
}

func main() {

	
	// 1. Step flag and start
	subdomain, port, err := skipperflag.StartSkipper()

	if err != nil {
		return
	}

	wg := &sync.WaitGroup{}
	// 2. Stablish connection with localhost
	config.localhostUrl = "localhost: " + strconv.Itoa(port)
	config.proxyUrl = "localhost:9000"
	_, err = net.DialTimeout("tcp", config.localhostUrl, time.Second*1)
	if err != nil {
		constants.PrintWithColor(constants.Red, "failed to connect to your localhost app. Check that the port is correct")
		return
	}

	wg.Add(1)
	// start a continous ping to it
	go PingLocalhost(config.localhostUrl, wg)

	// 3. start connection with proxy

	proxyConn, err := net.Dial("tcp", config.proxyUrl)
	if err != nil {
		fmt.Println("error ", err)
		return
	}

	// 4. send a frame to check to request connection and listen for the acknowledge
	RequestFrame := frame.CreateFrame(1, constants.Control_TunnelRequest, 0, uint32(len(subdomain)))
	RequestBuffer := RequestFrame.Encode([]byte(subdomain))
	proxyConn.Write(RequestBuffer)
	StartReadLoop(proxyConn)

	// 5. if we are here, means that we are already connected and ready, so we sping up the reading loop ffor every request
	go StartReadLoopConstant(proxyConn)
	wg.Wait()

}




func StartReadLoop(conn net.Conn) error {
	timeError := conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if timeError != nil {
		fmt.Println("ERROR on dedlinee")
		return fmt.Errorf(timeError.Error())
	}

	fmt.Println("we are on the reading loop for acknowledge")
	frameType, _, _, payload, err := frame.ReadCompleteFrame(conn)
	if err != nil {
		fmt.Println("error reading the proxy packet", err)
		return fmt.Errorf("error reading loop or connection closed", err)
	}

	switch frameType {
	case constants.Control_TunnelAck:
		fmt.Println("GOOD WE RECEIVED ALL GOOD")
		return nil
	case constants.Control_TunnelError:
		fmt.Println(string(payload))
		return err
	}

	fmt.Println("porque llegamos aca?")
	return err
}





// this is the loop with the reactor pattern
func StartReadLoopConstant(conn net.Conn) error {
	for {
		fmt.Println("we are on the reading loop general")
		frameType, _, _, payload, err := frame.ReadCompleteFrame(conn)
		if err != nil {
			fmt.Println("error reading the proxy packet", err)
			return fmt.Errorf("error reading loop or connection closed", err)
		}

		switch frameType {
		case constants.ProxyRequestType:
			// deserilaize and decode all and redierct to a workerpool

		case constants.ProxyPing:
			// this is just pinging
			fmt.Println(string(payload))
			return err
		}
		fmt.Println("porque llegamos aca?")
		return err
	}
}

func PingLocalhost(localhost string, wg *sync.WaitGroup) error {
	for {
		_, err := net.DialTimeout("tcp", localhost, time.Second*1)
		if err != nil {
			break
		}
		fmt.Println("ping checked")
		time.Sleep(time.Second * 4)
	}

	wg.Done()
	return fmt.Errorf("failed to conenct to your localhost app")
}
