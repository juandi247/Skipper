package http

import (
	"SkipperProxy/tunnel"
	"fmt"
	"net/http"
)

type httpServer struct {
	muxer     *http.ServeMux
	port      string
	startFunc func() error
	// RequestTimeout time

}

func CreateHttpServer(port string, tm tunnel.TunnelManager, isProd bool, certfile, keyfile string) *httpServer {
	httpMultiplexer := http.NewServeMux()
	// we register it as wildcard
	httpMultiplexer.HandleFunc("/", ClosureFunc(tm))
	fmt.Println("creamos server")

	server := &httpServer{
		muxer: httpMultiplexer,
		port:  port,
	}
	fmt.Println(isProd)

	if isProd {
		server.startFunc = func() error {
			fmt.Println("starting with tls")
			return http.ListenAndServeTLS(port, certfile, keyfile, httpMultiplexer)
		}
	} else {
		server.startFunc = func() error {
			fmt.Println("starting without tls, localhost")
			return http.ListenAndServe(port, httpMultiplexer)
		}
	}
	return server
}

func (s *httpServer) StartServer() error {
	err := s.startFunc()
	if err != nil {
		fmt.Println("ERORRRR", err)
		return fmt.Errorf(err.Error())
	}
	fmt.Println("empezo el de http")
	return nil
}
