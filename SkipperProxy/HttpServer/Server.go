package HttpServer

import (
	// tcpserver "SkipperTunnelProxy/TcpServer"
	"SkipperTunnelProxy/connectionmanager"
	"context"
	"fmt"
	"net/http"
	"html/template"
)

type Server struct {
	port   int
	Router *Router
	// middlewares  []Middleware
	errorHandler      ErrorHandler
	server            *http.Server
	certFile          string
	keyFile           string
	useHTTPS          bool
	ConnectionManager *connectionmanager.ConnectionManager
	Templates          *template.Template 

}

type ServerOption func(*Server)

func NewServer(port int, useHTTPS bool, connectionManager *connectionmanager.ConnectionManager) *Server {
	s := &Server{
		port: port,
		// TcpRequestChannel: ch,
		useHTTPS:          useHTTPS,
		ConnectionManager: connectionManager,
	}
	if useHTTPS {
		s.certFile = "/etc/letsencrypt/live/skipper.lat-0001/fullchain.pem"
		s.keyFile = "/etc/letsencrypt/live/skipper.lat-0001/privkey.pem"
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
