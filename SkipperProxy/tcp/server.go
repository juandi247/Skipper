package tcp

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	"SkipperProxy/tunnel"
	"fmt"
	"io"
	"net"
)

type TcpServer struct {
	port     string
	listener net.Listener
	cm 		*tunnel.ConnectionManager
}

func CreateTcpServer(port string, cm *tunnel.ConnectionManager) *TcpServer {
	return &TcpServer{
		port: port,
		cm: cm,
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

func (s *TcpServer)handleNewConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("new connectoin accepted, before handshake", conn)
	_, payload, err := ReadFrame(conn, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = s.cm.RegisterTunnel(string(payload), conn, conn.RemoteAddr())
	if err != nil {
		fmt.Println("ERROR regitrndo tunnel", err)
	}

	fmt.Println("tunnel registered Succesffulty and with a Acknoledge succesffull")

	ReadLoop(conn)
}

func ReadLoop(conn net.Conn) {
	for {
		f, payload, err := ReadFrame(conn, false)

		if err != nil {
			fmt.Println("error on reading loop", err)
			return
		}

		switch f.FrameType {
		case constants.TunnelResponseType:
			fmt.Println("llego la respuesta", payload)
			// send to the channel by processing the stream id

		case constants.TunnelPong:
			// todo: later
		}

	}
}

func ReadFrame(conn net.Conn, isControlType bool) (*frame.TcpFrame, frame.Payload, error) {

	buffer := make([]byte, 20)
	_, err := io.ReadFull(conn, buffer)

	if err != nil {
		return nil, nil, fmt.Errorf("ERROR Reading the buffer")
	}
	f, err := frame.DecodeHeader(buffer)

	if err != nil {
		return nil, nil, fmt.Errorf("could not parse the header", err)

	}
	err = constants.ValidateFramingType(f.FrameType)

	if err != nil {
		return nil, nil, fmt.Errorf("didnt receive right after the connection the REquest Type good")
	}

	payloadBuffer := make([]byte, f.PayloadLen)
	_, err = io.ReadFull(conn, payloadBuffer)
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR, the payload length was sent incomplete")

	}

	return f, payloadBuffer, nil
}
