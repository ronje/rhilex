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

package archsupport

type DeviceNodeType string

const (
	UART   DeviceNodeType = "UART"
	ETHNET DeviceNodeType = "ETHNET"
	NM4G   DeviceNodeType = "NM4G"
	NM5G   DeviceNodeType = "NM5G"
	CAN    DeviceNodeType = "CAN"
	RS485  DeviceNodeType = "RS485"
	RS232  DeviceNodeType = "RS232"
)

type DeviceNode struct {
	Name   string         `json:"name"`
	Type   DeviceNodeType `json:"type"`
	Status int            `json:"status"`
}

type DeviceTree struct {
	Network    []DeviceNode `json:"network"`
	Wlan       []DeviceNode `json:"wlan"`
	MNet4g     []DeviceNode `json:"net4g"`
	MNet5g     []DeviceNode `json:"net5g"`
	SerialPort []DeviceNode `json:"serialPort"`
	SoftRouter []DeviceNode `json:"soft_router"`
	CanBus     []DeviceNode `json:"canbus"`
}

func DefaultDeviceTree() DeviceTree {
	return DeviceTree{}
}
