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
)

func handleServer(server net.Conn, local net.Conn) {
	defer server.Close()
	defer local.Close()

	// 从服务端读取数据并发送到本地连接
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := server.Read(buffer)
			if n > 0 {
				_, err := local.Write(buffer[:n])
				if err != nil {
					log.Println("Error writing to local:", err)
					return
				}
			}
			if err != nil {
				log.Println("Error reading from server:", err)
				return
			}
		}
	}()

	// 从本地连接读取数据并发送到服务端
	buffer := make([]byte, 4096)
	for {
		n, err := local.Read(buffer)
		if n > 0 {
			_, err := server.Write(buffer[:n])
			if err != nil {
				log.Println("Error writing to server:", err)
				return
			}
		}
		if err != nil {
			log.Println("Error reading from local:", err)
			return
		}
	}
}

func StartClient() {
	// 连接到服务端
	server, err := net.Dial("tcp", "server_ip:8080")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer server.Close()

	// 进行认证
	if !authenticate(server) {
		log.Fatal("Authentication failed, exiting...")
	}

	// 要转发的本地地址
	localAddr := "127.0.0.1:80"

	// 发送本地地址到服务端
	_, err = server.Write([]byte(localAddr))
	if err != nil {
		log.Fatal("Error sending local address:", err)
	}

	// 连接到本地服务
	local, err := net.Dial("tcp", localAddr)
	if err != nil {
		log.Fatal("Error connecting to local service:", err)
	}
	defer local.Close()

	// 处理数据转发
	handleServer(server, local)
}
