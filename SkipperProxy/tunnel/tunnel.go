package tunnel

import "net"

type Tunnel struct {
	Subdomain  string
	Connection net.Conn
	ip         net.Addr
}

func CreateTunnel(subdomain string, conn net.Conn, ip net.Addr) *Tunnel {
	return &Tunnel{
		Subdomain:  subdomain,
		Connection: conn,
		ip:         ip,
	}
}
