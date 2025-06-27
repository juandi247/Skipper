package http

import (
	"fmt"
	"net/http"
)

type httpServer struct {
	muxer *http.ServeMux
	port  string
	// RequestTimeout time

}

func CreateHttpServer(port string) *httpServer {
	httpMultiplexer := http.NewServeMux()
	// we register it as wildcard
	httpMultiplexer.HandleFunc("/", handleRequest)
	fmt.Println("creamos server")
	return &httpServer{
		muxer: httpMultiplexer,
		port:  port,
	}
}

func (s *httpServer) StartServer() error {
	fmt.Println("intentmos empezarlo")
	err := http.ListenAndServe(s.port, s.muxer)
	if err != nil {
		fmt.Println("ERORRRR")
		return fmt.Errorf(err.Error())
	}
	fmt.Println("empezo el de http")
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ParseSubdomain(r.Host)
	fmt.Printf("the path of the request was %v \n", r.RequestURI)
	hello := "hello world"
	w.Write([]byte(hello))
}
