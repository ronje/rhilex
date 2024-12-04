// Copyright (C) 2024 wwhai
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

package test

import (
	"net"
	"testing"
	"time"

	"github.com/hootrhino/rhilex/component/tinydhcp"
)

// go test -timeout 30s -run ^TestTinyDhcpServer github.com/hootrhino/rhilex/test -v -count=1
func TestTinyDhcpServer(t *testing.T) {
	dhcpServer := &tinydhcp.DHCPServer{
		IPPools:        []*tinydhcp.IPPool{},
		StaticBindings: make(map[string]tinydhcp.StaticBinding),
		Log:            []string{},
		Gateway:        net.ParseIP("192.168.1.1"),
	}

	dhcpServer.AddIPPool("192.168.1.0/24", "192.168.1.100", "192.168.1.200", 24*time.Hour)
	mac, _ := net.ParseMAC("AA:BB:CC:DD:EE:01")
	dhcpServer.AddStaticBinding(mac, net.ParseIP("192.168.1.150"))
	dhcpServer.StartServer("eth1", 67)

}
