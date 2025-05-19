package connectionmanager

import (
	"SkipperTunnelProxy/message"
	"fmt"
	"net"
	"sync"
)

type ConnectionManager struct {
	TunnelConnectionsMap  map[string]net.Conn
	Mu                    sync.Mutex
	GlobalResponseChannel map[string]chan []byte
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		TunnelConnectionsMap:  make(map[string]net.Conn),
		Mu:                    *&sync.Mutex{},
		GlobalResponseChannel: make(map[string]chan []byte),
	}
}

func (cm *ConnectionManager) AddTunnelConnection(subdomain string, conn net.Conn) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	fmt.Println("FAMOS A AGREGAR")
	cm.TunnelConnectionsMap[subdomain] = conn
}

func (cm *ConnectionManager) DeleteTunnelConnection(subdomain string) {
	delete(cm.TunnelConnectionsMap, subdomain)
}

func (cm *ConnectionManager) SendMessageToTunnel(subdomain string, message message.TcpMessage) error {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	conn, _ := cm.TunnelConnectionsMap[subdomain]
	fmt.Println("SUBDOMINOOOO", subdomain)
	fmt.Println(message.Data, "MESAGE DATAAAAAAAAAAA")
	_, err := conn.Write(message.Data)
	if err != nil {
		return fmt.Errorf("Error escribiendo a conexión TCP: %v\n", err)
	}
	return nil
}

func (cm *ConnectionManager) GetTunnelConnection(subdomain string) (bool, error) {
	_, ok := cm.TunnelConnectionsMap[subdomain]
	fmt.Println("MAPAPAA")

	fmt.Println(&cm.TunnelConnectionsMap)
	fmt.Println(ok, "aa")

	if !ok {
		return false, fmt.Errorf("subdomain doesnt exist")
	}
	fmt.Println("AAAAA", cm.TunnelConnectionsMap)

	return true, nil

}

func (cm *ConnectionManager) SaveResponseChannel(uid string, ch chan []byte) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.GlobalResponseChannel[uid] = ch
}
