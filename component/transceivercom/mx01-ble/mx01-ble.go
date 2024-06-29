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
	serial "github.com/hootrhino/goserial"
	mx01 "github.com/hootrhino/rhilex-goat/bsp/mx01"
	"github.com/hootrhino/rhilex-goat/device"
	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"os"
	"time"
)

type Mx01BLEConfig struct {
	ComConfig transceivercom.TransceiverConfig
}
type Mx01BLE struct {
	R          typex.Rhilex
	mainConfig Mx01BLEConfig
	mx01       device.Device
}

func NewMx01BLE(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &Mx01BLE{R: R, mainConfig: Mx01BLEConfig{
		ComConfig: transceivercom.TransceiverConfig{
			Address:   "COM3",
			BaudRate:  9600,
			DataBits:  8,
			Parity:    "N",
			StopBits:  1,
			IOTimeout: 50,  // IOTimeout * time.Millisecond
			ATTimeout: 200, // ATRwTimeout * time.Millisecond
		},
	}}
}
func (tc *Mx01BLE) Start(transceivercom.TransceiverConfig) error {
	env := os.Getenv("BLESUPPORT")
	if env == "MX01" {
		glogger.GLogger.Info("MX01-BLE-Module Init")
		config := serial.Config{
			Address:  tc.mainConfig.ComConfig.Address,
			BaudRate: tc.mainConfig.ComConfig.BaudRate,
			DataBits: tc.mainConfig.ComConfig.DataBits,
			Parity:   tc.mainConfig.ComConfig.Parity,
			StopBits: tc.mainConfig.ComConfig.StopBits,
			Timeout:  time.Duration(tc.mainConfig.ComConfig.IOTimeout) * time.Millisecond,
		}
		serialPort, err := serial.Open(&config)
		if err != nil {
			return err
		}
		tc.mx01 = mx01.NewMX01("mx01", serialPort)
		tc.mx01.Flush()
		glogger.GLogger.Info("MX01-BLE-Module Init Ok.")
	}
	glogger.GLogger.Info("MX01-BLE-Module Started")
	return nil
}
func (tc *Mx01BLE) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (tc *Mx01BLE) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "MX01-BLE-Module",
		Model:    "MX-01",
		Type:     transceivercom.BLE,
		Vendor:   "SHENZHEN-MIAOXIANG-technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *Mx01BLE) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_ERROR,
		Error: nil,
	}
}
func (tc *Mx01BLE) Stop() {
	glogger.GLogger.Info("MX01-BLE-Module Stopped")
	if tc.mx01 != nil {
		tc.mx01.Close()
	}
}
