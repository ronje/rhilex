// Copyright (C) 2025 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package microtun

import (
	"log"
	"net"
	"sync"
)

// Client 表示一个连接的客户端
type Client struct {
	Conn    net.Conn
	Target  string
	IsAlive bool
}

// ClientManager 管理所有连接的客户端
type ClientManager struct {
	Clients map[string]*Client
	Mutex   sync.Mutex
}

// AddClient 添加一个新的客户端
func (cm *ClientManager) AddClient(id string, client *Client) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Clients[id] = client
}

// RemoveClient 移除一个客户端
func (cm *ClientManager) RemoveClient(id string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	if client, ok := cm.Clients[id]; ok {
		client.Conn.Close()
		delete(cm.Clients, id)
	}
}

// GetClient 获取一个客户端
func (cm *ClientManager) GetClient(id string) (*Client, bool) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	client, ok := cm.Clients[id]
	return client, ok
}

var clientManager = ClientManager{
	Clients: make(map[string]*Client),
}

func handleClient(client net.Conn, target net.Conn, clientID string) {
	defer client.Close()
	defer target.Close()
	defer clientManager.RemoveClient(clientID)

	// 从客户端读取数据并发送到目标连接
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := client.Read(buffer)
			if n > 0 {
				_, err := target.Write(buffer[:n])
				if err != nil {
					log.Println("Error writing to target:", err)
					return
				}
			}
			if err != nil {
				log.Println("Error reading from client:", err)
				return
			}
		}
	}()

	// 从目标连接读取数据并发送到客户端
	buffer := make([]byte, 4096)
	for {
		n, err := target.Read(buffer)
		if n > 0 {
			_, err := client.Write(buffer[:n])
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		}
		if err != nil {
			log.Println("Error reading from target:", err)
			return
		}
	}
}

func StartServer() {
	// 监听服务端端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	log.Println("Server is listening on :8080")

	for {
		// 接受客户端连接
		client, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// 进行认证
		if !authenticate(client) {
			client.Close()
			continue
		}

		// 等待客户端发送目标地址
		buffer := make([]byte, 1024)
		n, err := client.Read(buffer)
		if err != nil {
			log.Println("Error reading target address:", err)
			client.Close()
			continue
		}
		targetAddr := string(buffer[:n])

		// 连接到目标地址
		target, err := net.Dial("tcp", targetAddr)
		if err != nil {
			log.Println("Error connecting to target:", err)
			client.Close()
			continue
		}

		// 生成客户端ID
		clientID := client.RemoteAddr().String()

		// 添加客户端到管理列表
		clientManager.AddClient(clientID, &Client{
			Conn:    client,
			Target:  targetAddr,
			IsAlive: true,
		})

		// 处理数据转发
		go handleClient(client, target, clientID)
	}
}
