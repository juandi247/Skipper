/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"SkipperTunnel/cmd"
)

var Env string

func main() {
	cmd.Env= Env
	cmd.Execute()
}
