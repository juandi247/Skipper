package tunnel

import (
	"context"
	"net"
)

type Tunnel struct {
	localhostUrl string
	proxyConn    net.Conn
	ctx          context.Context
}

func NewTunnel(localhost string, proxyConn net.Conn, ctx context.Context) *Tunnel {
	return &Tunnel{
		localhostUrl: localhost,
		proxyConn:    proxyConn,
		ctx:          ctx,
	}
}
