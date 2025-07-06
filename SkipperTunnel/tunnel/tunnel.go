package tunnel

import (
	"context"
	"net"
)

type Tunnel struct {
	Subdomain string
	LocalhostUrl string
	ProxyUrl string
	ProxyConn    net.Conn
	Ctx          context.Context
}

func NewTunnel(proxyUrl string, ctx context.Context) *Tunnel {
	return &Tunnel{
		ProxyUrl: proxyUrl,
		Ctx:          ctx,
	}
}
