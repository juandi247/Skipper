/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"SkipperTunnel/tunnel"
	"context"
	"fmt"
)

func main() {

	ctx:= context.Background()
	ctx, cancel:= context.WithCancel(ctx)
	errChan:= make(chan error)

	tn:= tunnel.NewTunnel("skipper.lat",ctx)

	tn.FsmStart()
	// todo: add the buffered chanel for errors that will tr


	errorRecevied:= <- errChan
	cancel()
	fmt.Println("cancelling all the thigns")
	}
}