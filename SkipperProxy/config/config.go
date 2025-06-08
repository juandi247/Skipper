package config

import "fmt"

type Config struct {
	TcpPort      string
	HttpPort     int16
	BaseDomain   string
	IsHttps      bool
	WorkerNumber int
	DomainParts  int
}

var devConfig Config = Config{
	TcpPort:      ":9000",
	HttpPort:     8080,
	BaseDomain:   "localhost:8080",
	IsHttps:      false,
	WorkerNumber: 20,
	DomainParts:  1,
}

var prodConfig Config = Config{
	TcpPort:      ":9000",
	HttpPort:     443,
	BaseDomain:   "skipper.lat",
	IsHttps:      true,
	WorkerNumber: 100,
	DomainParts:  2,
}

func LoadConfig(environment string) Config {
	var config Config
	switch environment {
	case "DEV":
		fmt.Println("we are on dev mode")
		config = devConfig
	case "PROD":
		fmt.Println("we are on prod")
		config = prodConfig
	default:
		config = devConfig
	}
	return config
}
