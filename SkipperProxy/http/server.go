package http

import (
	"SkipperProxy/tunnel"
	"fmt"
	"net/http"
)

type httpServer struct {
	muxer *http.ServeMux
	port  string
	// RequestTimeout time

}

func CreateHttpServer(port string, tm tunnel.TunnelManager) *httpServer {
	httpMultiplexer := http.NewServeMux()
	// we register it as wildcard
	httpMultiplexer.HandleFunc("/", ClosureFunc(tm))
	fmt.Println("creamos server")
	return &httpServer{
		muxer: httpMultiplexer,
		port:  port,
	}
}

func (s *httpServer) StartServer() error {
	err := http.ListenAndServe(s.port, s.muxer)
	if err != nil {
		fmt.Println("ERORRRR", err)
		return fmt.Errorf(err.Error())
	}
	fmt.Println("empezo el de http")
	return nil
}
