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
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	mx01 "github.com/hootrhino/rhilex-goat/bsp/mx01"
	"github.com/hootrhino/rhilex-goat/device"
	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type Mx01BLEConfig struct {
	ComConfig transceivercom.TransceiverConfig
}
type Mx01BLE struct {
	R          typex.Rhilex
	mainConfig Mx01BLEConfig
	mx01       device.Device
	locker     sync.Mutex
}

func NewMx01BLE(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &Mx01BLE{R: R, locker: sync.Mutex{}, mainConfig: Mx01BLEConfig{
		ComConfig: transceivercom.TransceiverConfig{
			Address:   "COM1",
			BaudRate:  9600,
			DataBits:  8,
			Parity:    "N",
			StopBits:  1,
			IOTimeout: 50,  // IOTimeout * time.Millisecond
			ATTimeout: 200, // ATRwTimeout * time.Millisecond
		},
	}}
}
func (tc *Mx01BLE) Start(Config transceivercom.TransceiverConfig) error {
	env := os.Getenv("BLESUPPORT")
	if env == "MX01" {
		glogger.GLogger.Info("MX01-BLE Init")
		config := serial.Config{
			Address:  Config.Address,
			BaudRate: Config.BaudRate,
			DataBits: Config.DataBits,
			Parity:   Config.Parity,
			StopBits: Config.StopBits,
			Timeout:  time.Duration(tc.mainConfig.ComConfig.IOTimeout) * time.Millisecond,
		}
		serialPort, err := serial.Open(&config)
		if err != nil {
			return err
		}
		tc.mx01 = mx01.NewMX01("mx01", serialPort)
		tc.mx01.Flush()
		go func(io io.ReadWriteCloser) {
			for {
				N, Bytes := utils.ReadInLeastTimeout(context.Background(), io,
					time.Duration(tc.mainConfig.ComConfig.ATTimeout)*time.Millisecond)
				if N > 0 {
					glogger.GLogger.Debug("MX01-BLE Read Data: ", Bytes[:N])
				}
			}

		}(serialPort)
		glogger.GLogger.Info("MX01-BLE Started")
	}
	return nil
}
func (tc *Mx01BLE) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (tc *Mx01BLE) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "mx01",
		Model:    "MX01-BLE",
		Type:     transceivercom.BLE,
		Vendor:   "SHENZHEN-MIAOXIANG-technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *Mx01BLE) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_ERROR,
		Error: fmt.Errorf("NOT SUPPORT"),
	}
}
func (tc *Mx01BLE) Stop() {
	glogger.GLogger.Info("MX01-BLE Stopped")
	if tc.mx01 != nil {
		tc.mx01.Close()
	}
}
