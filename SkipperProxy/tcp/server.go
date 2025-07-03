package tcp

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	"SkipperProxy/tunnel"
	"fmt"
	"net"
)

type Server interface {
	StartServer() error
}

type TcpServer struct {
	port     string
	listener net.Listener
	tm       tunnel.TunnelManager
}

func CreateTcpServer(port string, tm tunnel.TunnelManager) *TcpServer {
	return &TcpServer{
		port: port,
		tm:   tm,
	}
}

func (s *TcpServer) StartServer() {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		panic("Tcp server could not start!")
	}
	s.listener = ln

	s.AcceptLoop()

}

func (s *TcpServer) AcceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleNewConnection(conn)
	}

}

func (s *TcpServer) handleNewConnection(conn net.Conn) {
	defer func() {
		fmt.Println("vamos a cerrar la conexion", conn)
		conn.Close()
	}()
	frameType, _, _, payload, err := frame.ReadCompleteFrame(conn, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	// assertion for the tunnel (it should ALWAYS send that first)
	if frameType != constants.Control_TunnelRequest {
		panic("the tunnel didnt sent the request connection type")
	}

	tunnelConnection, err := s.tm.RegisterTunnel(string(payload), conn, conn.RemoteAddr())
	if err != nil {
		fmt.Println("ERROR regitrndo tunnel", err)
	}

	fmt.Println("tunnel registered Succesffulty and with a Acknoledge succesffull")

	tunnelConnection.StartReadLoop()
}
