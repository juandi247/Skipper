package tcpserver

import (
	"fmt"
	"net"
	// "time"

	// "sync"
	"SkipperTunnelProxy/connectionmanager"
)

// todo: change the type of the channel isntead of bytes, use the json parsing but later

type Server struct {
	listenAdrr string
	ln         net.Listener
	quitch     chan struct{}
	//good practice to use bytes
	MessageChanel chan []byte
	// todo: a connection map to help later with the wildcard and redirecting routing
	ConnectionManager *connectionmanager.ConnectionManager
	// ConnMutext     sync.Mutex
	// ConnectionMap  map[string]net.Conn
	// RequestChannel chan TcpMessage
}

func NewServer(listenAddr string, cm *connectionmanager.ConnectionManager) *Server {
	return &Server{
		listenAdrr:        listenAddr,
		quitch:            make(chan struct{}),
		MessageChanel:     make(chan []byte, 10),
		ConnectionManager: cm,
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
		addr := conn.RemoteAddr().String()
		s.ConnectionManager.AddTunnelConnection(addr, conn)
		// accept the subdomain
		// Procesar el subdominio antes de comenzar el loop de lectura
		subdomain, err := s.ReviewTunnelConnection(conn)
		if err != nil {
			fmt.Println("Failed to add subdomain:", err)
			conn.Close()
			continue
		}

		fmt.Println("Successfully added subdomain:", subdomain)
		s.ConnectionManager.AddTunnelConnection(subdomain, conn)

		fmt.Println("no errros ading subdomain an now listning with the read loop")
		go s.ReadLoop(conn, subdomain)
	}
}

func (s *Server) ReviewTunnelConnection(conn net.Conn) (string, error) {
	buffer := make([]byte, 2048)
	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Leer el subdominio
	numberOfBytes, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("ERROR LEY_EDNO")
		return "", fmt.Errorf("error reading subdomain: %v", err)
	}

	subdomain := string(buffer[:numberOfBytes])

	// Verificar si el subdominio ya está en uso
	if _, exists := s.ConnectionManager.TunnelConnectionsMap[subdomain]; exists {
		conn.Write([]byte("The subdomain is already in use"))
		return "", fmt.Errorf("subdomain already used: %s", subdomain)
	}

	fmt.Println("YA TERMINO REVIEW TODO BIEN")

	return subdomain, nil
}

func (s *Server) ReadLoop(conn net.Conn, subdomain string) {
	defer conn.Close()

	defer func() {
		conn.Close()
		addr := conn.RemoteAddr().String()
		fmt.Println("Connection closed:", addr)
		s.ConnectionManager.DeleteTunnelConnection(subdomain)

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
