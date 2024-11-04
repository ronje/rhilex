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

package dnsserver

import (
	"fmt"
	"os"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"gopkg.in/ini.v1"
)

type TinyDnsServer struct {
	server *Server
}

func NewTinyDnsServer() *TinyDnsServer {
	return &TinyDnsServer{}
}

func (dm *TinyDnsServer) Init(config *ini.Section) error {
	return nil
}

func (dm *TinyDnsServer) Start(typex.Rhilex) error {
	host, _ := os.Hostname()
	info := []string{"rhilex-service"}
	service, _ := NewMDNSService(host, fmt.Sprintf("rhilex.service.%s", host), "", "", 40000, nil, info)
	dm.server, _ = NewServer(&Config{Zone: service, Logger: glogger.GLogger.Writer()})
	return nil
}
func (dm *TinyDnsServer) Stop() error {
	if dm.server != nil {
		dm.server.Shutdown()
	}
	return nil
}

func (dm *TinyDnsServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "DNS-SERVER",
		Name:        "TinyDnsServer",
		Version:     "v0.0.1",
		Description: "Tiny Dns Server",
	}
}

/*
*
* 服务调用接口
*
 */
func (dm *TinyDnsServer) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
