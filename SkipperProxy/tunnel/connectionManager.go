package tunnel

import (
	"SkipperProxy/frame"
	"fmt"
	"net"
	"sync"
)

type ConnectionManager struct {
	TunnelConnectionsMap map[string]net.Conn
	Mutex                sync.Mutex
}

func CreateConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		TunnelConnectionsMap: make(map[string]net.Conn, 100),
		Mutex:                sync.Mutex{},
	}
}

func (cm *ConnectionManager) RegisterTunnel(subdomain string, conn net.Conn, ip net.Addr) (*Tunnel, error) {
	if len(subdomain) == 0 {
		return nil, fmt.Errorf("subdomain cannot be empty")
	}
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	if _, exists := cm.TunnelConnectionsMap[subdomain]; exists {
		frame.SendControlConnectionError(conn, "The subdomain is already taken, please try again")
		conn.Close()
		return nil, fmt.Errorf("subdomain already in use")
	}

	tunnel := CreateTunnel(subdomain, conn, ip)
	cm.TunnelConnectionsMap[subdomain] = conn
	frame.SendControlOk(conn)
	return tunnel, nil
}

func (cm *ConnectionManager) RemoveTunnel(tunnel *Tunnel) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	defer tunnel.Connection.Close()
	delete(cm.TunnelConnectionsMap, tunnel.Subdomain)
}
