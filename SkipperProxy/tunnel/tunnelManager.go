package tunnel

import (
	"fmt"
	"net"
	"sync"
)

type TunnelManager struct {
	TunnelConnectionsMap map[string]*TunnelConnection
	Mutex                sync.Mutex
}

func CreateTunnelManager() *TunnelManager {
	return &TunnelManager{
		TunnelConnectionsMap: make(map[string]*TunnelConnection, 100),
		Mutex:                sync.Mutex{},
	}
}

func (tm *TunnelManager) RegisterTunnel(subdomain string, conn net.Conn, ip net.Addr) (*TunnelConnection, error) {
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

func (tm *TunnelManager) RemoveTunnel(tunnel *TunnelConnection) {
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()
	defer tunnel.Connection.Close()
	delete(tm.TunnelConnectionsMap, tunnel.Subdomain)
}
