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

package mx01ble

import (
	"time"

	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/typex"
)

type Mx01BLEConfig struct {
}
type Mx01BLE struct {
	R          typex.Rhilex
	mainConfig Mx01BLEConfig
}

func NewMx01BLE(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &Mx01BLE{R: R, mainConfig: Mx01BLEConfig{}}
}
func (tc *Mx01BLE) Start(transceivercom.TransceiverConfig) error {
	return nil
}
func (tc *Mx01BLE) Ctrl(cmd []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (tc *Mx01BLE) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:   "MX01-BLE-Module",
		Model:  "MX-01S",
		Type:   transceivercom.BLE,
		Vendor: "SHENZHEN-MIAOXIANG-TECH",
	}
}
func (tc *Mx01BLE) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_UP,
		Error: nil,
	}
}
func (tc *Mx01BLE) Stop() {

}
