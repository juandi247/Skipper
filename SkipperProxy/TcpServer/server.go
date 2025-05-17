package tcpserver

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// todo: change the type of the channel isntead of bytes, use the json parsing but later

type Server struct {
	listenAdrr string
	ln         net.Listener
	quitch     chan struct{}
	//good practice to use bytes
	MessageChanel chan []byte
	// todo: a connection map to help later with the wildcard and redirecting routing
	ConnMutext     sync.Mutex
	ConnectionMap  map[string]net.Conn
	RequestChannel chan TcpMessage
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAdrr:     listenAddr,
		quitch:         make(chan struct{}),
		MessageChanel:  make(chan []byte, 10),
		ConnectionMap:  make(map[string]net.Conn),
		RequestChannel: make(chan TcpMessage),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAdrr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.AcceptLoop()

	<-s.quitch

	// close the channel when we stop the server
	close(s.quitch)
	return nil
}

func (s *Server) AcceptLoop() {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}
		fmt.Println("New connection", conn.RemoteAddr())
		// addr := conn.RemoteAddr().String()
		buffer := make([]byte, 2048)
		numberOfBytes, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("read error", err)
			conn.Close()
			continue
		}

		// Convertir el buffer a string para extraer el subdominio
		subdomain := strings.TrimSpace(string(buffer[:numberOfBytes]))
		fmt.Println("Subdominio recibido:", subdomain)


		s.ConnMutext.Lock()
		if _, exists := s.ConnectionMap[subdomain]; exists {
			fmt.Println("⚠️ Error: Subdominio ya en uso:", subdomain)
			s.ConnMutext.Unlock()
			conn.Close()
			continue
		}
		s.ConnectionMap[subdomain] = conn
		s.ConnMutext.Unlock()
		fmt.Println(subdomain)
		// handling every connection on a different goroutine
		go s.ReadLoop(conn)
	}

}

func (s *Server) ReadLoop(conn net.Conn) {
	defer func() {
		conn.Close()
		addr := conn.RemoteAddr().String()
		s.ConnMutext.Lock()
		delete(s.ConnectionMap, addr)
		s.ConnMutext.Unlock()
		fmt.Println("Connection closed:", addr)
	}()

	buffer := make([]byte, 2048)
	for {
		numberOfBytes, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("read error", err)
			break
		}
		// write to the tcp client
		// conn.Write([]byte(string("te escribo de vuelta!")))

		// send the message received to the channel
		s.MessageChanel <- buffer[:numberOfBytes]
	}

}
