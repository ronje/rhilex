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
	"testing"

	"github.com/hootrhino/rhilex/periphery/rhilexg1"
)

// go test -timeout 30s -run ^Test_gen_eth_config github.com/hootrhino/rhilex/test -v -count=1
func Test_gen_eth_config(t *testing.T) {

	configs := []rhilexg1.NetworkInterfaceConfig{
		{
			Interface:   "eth0",
			Address:     "192.168.1.100",
			Netmask:     "255.255.255.0",
			Gateway:     "192.168.1.1",
			DHCPEnabled: false,
		},
		{
			Interface:   "eth1",
			Address:     "192.168.1.101",
			Netmask:     "255.255.255.0",
			Gateway:     "192.168.1.1",
			DHCPEnabled: false,
		},
		{
			Interface:   "eth2",
			Address:     "192.168.1.101",
			Netmask:     "255.255.255.0",
			Gateway:     "192.168.1.1",
			DHCPEnabled: true,
		},
	}
	if err := rhilexg1.SetEthernet(configs); err != nil {
		t.Fatal(err.Error())
	}
}
