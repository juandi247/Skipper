/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"SkipperTunnel/frame"
	"SkipperTunnel/tunnel"
	"context"
	"fmt"
	"sync"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error, 1)
	requestChan:= make(chan *frame.InternalFrame, 1)
	syncPool:= &sync.Pool{New: func() interface{} {
		return &frame.InternalFrame{}
	}}


	defer func() {
		close(errChan)
		close(requestChan)
		cancel()
		fmt.Println("sendet the CANCEL!")
	}()

	tn := tunnel.NewTunnel("localhost:9000", ctx, errChan, requestChan, syncPool)
	tn.FsmStart()
}
