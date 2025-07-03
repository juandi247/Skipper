package tunnel

import (
	"fmt"
	"net"
	"sync"
)


type TunnelManager interface {
    RegisterTunnel(subdomain string, conn net.Conn, ip net.Addr) (*TunnelConnection, error)
    RemoveTunnel(tunnel *TunnelConnection)
    GetTunnel(subdomain string) (*TunnelConnection, error)
}

type SkipperManager struct {
	TunnelConnectionsMap map[string]*TunnelConnection
	Mutex                sync.RWMutex
}

func CreateSkipperManager() *SkipperManager {
	return &SkipperManager{
		TunnelConnectionsMap: make(map[string]*TunnelConnection, 100),
		Mutex:                sync.RWMutex{},
	}
}

func (tm *SkipperManager) RegisterTunnel(subdomain string, conn net.Conn, ip net.Addr) (*TunnelConnection, error) {
	if len(subdomain) == 0 {
		return nil, fmt.Errorf("subdomain cannot be empty")
	}
	tm.Mutex.Lock()
	if _, exists := tm.TunnelConnectionsMap[subdomain]; exists {
		tm.Mutex.Unlock()
		// this method isnt part of the TunnelConnection because we should not create a tunnel connection without a valid subdomain and validations
		SendControlConnectionError(conn, "The subdomain is already taken, please try again")
		conn.Close()
		return nil, fmt.Errorf("subdomain already in use")
	}
	tunnelConnection := CreateTunnelConnection(subdomain, conn, ip)

	tm.TunnelConnectionsMap[subdomain] = tunnelConnection
	tm.Mutex.Unlock()

	err := tunnelConnection.SendAcknowledgeConnection()
	if err != nil {
		return nil, fmt.Errorf("error creating tunnelConnection", tunnelConnection)
	}
	return tunnelConnection, nil
}

func (tm *SkipperManager) RemoveTunnel(tunnel *TunnelConnection) {
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()
	defer tunnel.Connection.Close()
	delete(tm.TunnelConnectionsMap, tunnel.Subdomain)
}

func (tm *SkipperManager) GetTunnel(subdomain string) (*TunnelConnection, error) {
	tm.Mutex.Lock()
	tunnel, exists := tm.TunnelConnectionsMap[subdomain]
	tm.Mutex.Unlock()
	if !exists {
		return nil, fmt.Errorf("The subdomain doens exists on the tunnel connections map")
	}
	return tunnel, nil
}
