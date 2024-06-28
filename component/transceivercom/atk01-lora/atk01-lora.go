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

package atk01lora

import (
	"time"

	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

type ATK01LoraConfig struct {
}
type ATK01Lora struct {
	R          typex.Rhilex
	mainConfig ATK01LoraConfig
}

func NewATK01Lora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &ATK01Lora{R: R, mainConfig: ATK01LoraConfig{}}
}
func (tc *ATK01Lora) Start(transceivercom.TransceiverConfig) error {
	glogger.GLogger.Info("EC200ADtu Started")
	return nil
}
func (tc *ATK01Lora) Ctrl(cmd []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (tc *ATK01Lora) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "ATK-01-LORA",
		Model:    "ATK-01-SX1278",
		Type:     transceivercom.LORA,
		Vendor:   "GUANGZHOU-ZHENGDIAN-YUANZI technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *ATK01Lora) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_UP,
		Error: nil,
	}
}
func (tc *ATK01Lora) Stop() {
	glogger.GLogger.Info("EC200ADtu Stopped")
}
