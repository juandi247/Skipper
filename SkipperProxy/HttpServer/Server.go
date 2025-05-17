package HttpServer

import (
	tcpserver "SkipperTunnelProxy/TcpServer"
	"context"
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	port   int
	Router *Router
	// middlewares  []Middleware
	errorHandler      ErrorHandler
	server            *http.Server
	TcpRequestChannel chan tcpserver.TcpMessage
	certFile          string
	keyFile           string
	useHTTPS          bool
}

var responseMap = make(map[string]chan []byte)
var responseMapMutex sync.Mutex

type ServerOption func(*Server)

func NewServer(port int, ch chan tcpserver.TcpMessage, useHTTPS bool) *Server {
	s := &Server{
		port:              port,
		TcpRequestChannel: ch,
		useHTTPS:          useHTTPS,
	}
	if useHTTPS {
		s.certFile = "/etc/letsencrypt/live/skipper.lat/fullchain.pem" // Ruta del certificado
		s.keyFile = "/etc/letsencrypt/live/skipper.lat/privkey.pem"    // Ruta de la clave privada
	}
	s.Router = NewRouter(s)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	/* check this server from http, because the doc says that if we leave it blanc,
	   takes port 80
	*/
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}
	fmt.Printf("Server starting on port %d", s.port)

	if s.useHTTPS {
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	return s.server.ListenAndServe()

}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) GetRouter() *Router {
	return s.Router
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func withErrorHandler(handler ErrorHandler) ServerOption {
	return func(s *Server) {
		s.errorHandler = handler
	}
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) RegisterResponseChannel(requestID string, ch chan []byte) {
	responseMapMutex.Lock()
	defer responseMapMutex.Unlock()
	responseMap[requestID] = ch
}

func (s *Server) GetResponseChannel(requestID string) (chan []byte, bool) {
	responseMapMutex.Lock()
	defer responseMapMutex.Unlock()
	ch, ok := responseMap[requestID]
	if ok {
		delete(responseMap, requestID)
	}
	return ch, ok
}
