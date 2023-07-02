package znet

import (
	"fmt"
	"sync"

	"github.com/pingtouge2008/zinx/ziface"
)

type ConnectionManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnectionManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully, conn num: ", cm.Len())
}

func (cm *ConnectionManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection remove from ConnManager successfully, conn num: ", cm.Len())
}

func (cm *ConnectionManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("connection %d is not exist", connID)
}

func (cm *ConnectionManager) Len() int {
	return len(cm.connections)
}
func (cm *ConnectionManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for _, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, conn.GetConnID())
	}

	fmt.Println("connection clear from ConnManager successfully, conn num: ", cm.Len())
}
