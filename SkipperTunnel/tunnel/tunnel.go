package tunnel

import (
	"SkipperTunnel/frame"
	"context"
	"net"
	"net/http"
	"sync"
)

type Tunnel struct {
	Subdomain            string
	LocalhostUrl         string
	ProxyUrl             string
	ProxyConn            net.Conn
	Ctx                  context.Context
	ErrChan              chan error
	RequestChan          chan *frame.InternalFrame
	syncPool             *sync.Pool
	DashboardRequestChan chan *http.Request
}

func NewTunnel(proxyUrl string, ctx context.Context, errChan chan error, requestChan chan *frame.InternalFrame, sp *sync.Pool) *Tunnel {
	return &Tunnel{
		ProxyUrl:    proxyUrl,
		Ctx:         ctx,
		ErrChan:     errChan,
		RequestChan: requestChan,
		syncPool:    sp,
	}
}
