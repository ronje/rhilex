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
	"bufio"
	"encoding/json"
	"log"
	"net"
)

// AuthRequest 认证请求结构体
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse 认证响应结构体
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func authenticate(client net.Conn) bool {
	reader := bufio.NewReader(client)
	// 读取客户端的认证请求
	data, err := reader.ReadBytes('\n')
	if err != nil {
		log.Println("Error reading authentication request:", err)
		return false
	}

	var authRequest AuthRequest
	err = json.Unmarshal(data, &authRequest)
	if err != nil {
		log.Println("Error unmarshaling authentication request:", err)
		return false
	}

	// 验证用户名和密码
	success := authRequest.Username == "username" && authRequest.Password == "password"
	response := AuthResponse{
		Success: success,
		Message: "",
	}
	if !success {
		response.Message = "Invalid username or password"
	}

	// 发送认证响应给客户端
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling authentication response:", err)
		return false
	}
	responseData = append(responseData, '\n')
	_, err = client.Write(responseData)
	if err != nil {
		log.Println("Error sending authentication response:", err)
		return false
	}

	return success
}
