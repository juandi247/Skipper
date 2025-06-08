package config

type Config struct {
	ProxyUrl string
	Workers  int
}

var devConfig Config = Config{
	ProxyUrl: "localhost:8080",
	Workers:  20,
}

var prodConfig Config = Config{
	ProxyUrl: "skipper.lat:9000",
	Workers:  100,
}

func LoadConfig(env string) Config {
	var config Config
	switch env {
	case "DEV":
		config = devConfig
	case "PROD":
		config = prodConfig
	default:
		config = devConfig
	}
	return config

}
