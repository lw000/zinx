package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/aceld/zinx/ziface"
)

//ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

//NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

//Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

//Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

// GetAllConnID 获取所有连接的ID
func (connMgr *ConnManager) GetAllConnID() []uint32 {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	ids := make([]uint32, 0, len(connMgr.connections))

	for id := range connMgr.connections {
		ids = append(ids, id)
	}

	return ids
}
