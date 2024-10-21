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
	"context"
	"fmt"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/component/internotify"
	transceivercom "github.com/hootrhino/rhilex/component/transceiver"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/protocol"
	"github.com/hootrhino/rhilex/typex"
)

type ATK01LoraConfig struct {
	ComConfig transceivercom.TransceiverConfig
}
type ATK01Lora struct {
	R              typex.Rhilex
	locker         sync.Mutex
	mainConfig     ATK01LoraConfig
	ProtocolSlaver *protocol.GenericProtocolSlaver
}

func NewATK01Lora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &ATK01Lora{
		R:      R,
		locker: sync.Mutex{},
		mainConfig: ATK01LoraConfig{
			ComConfig: transceivercom.TransceiverConfig{
				Address:           "COM1",
				BaudRate:          9600,
				DataBits:          8,
				Parity:            "N",
				StopBits:          1,
				IOTimeout:         50,  // IOTimeout * time.Millisecond
				ATTimeout:         200, // ATRwTimeout * time.Millisecond
				TransportProtocol: 1,
			},
		}}
}

/*
*
* Start
*
 */
func (tc *ATK01Lora) Start(Config transceivercom.TransceiverConfig) error {
	tc.mainConfig = ATK01LoraConfig{
		ComConfig: Config,
	}
	glogger.GLogger.Info("ATK01-LORA-01 Init")
	serialPort, errOpen := serial.Open(&serial.Config{
		Address:  tc.mainConfig.ComConfig.Address,
		BaudRate: tc.mainConfig.ComConfig.BaudRate,
		DataBits: tc.mainConfig.ComConfig.DataBits,
		Parity:   tc.mainConfig.ComConfig.Parity,
		StopBits: tc.mainConfig.ComConfig.StopBits,
		Timeout:  time.Duration(tc.mainConfig.ComConfig.IOTimeout) * time.Millisecond,
	})
	if errOpen != nil {
		return errOpen
	}
	config := protocol.TransporterConfig{
		Port:         serialPort,
		ReadTimeout:  time.Duration(tc.mainConfig.ComConfig.IOTimeout),
		WriteTimeout: time.Duration(tc.mainConfig.ComConfig.IOTimeout),
	}
	ctx, cancel := context.WithCancel(context.Background())
	tc.ProtocolSlaver = protocol.NewGenericProtocolSlaver(ctx, cancel, config)
	go tc.ProtocolSlaver.StartLoop(func(AppLayerFrame protocol.AppLayerFrame) {
		buffer, _ := AppLayerFrame.Encode()
		glogger.GLogger.Debug("ATK01Lora.ProtocolSlaver.Receive:", AppLayerFrame.String())
		internotify.Push(internotify.BaseEvent{
			Type:    "transceiver.up.data",
			Event:   "transceiver.up.data.atk01",
			Ts:      uint64(time.Now().UnixMilli()),
			Summary: "transceiver.up.data",
			Info:    buffer,
		})
	})
	glogger.GLogger.Info("ATK01-LORA-01 Started")
	return nil
}

func (tc *ATK01Lora) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	if string(topic) == "lora.atk01.cmd.send" {
		return []byte("OK"), nil
	}
	return []byte("OK"), nil
}
func (tc *ATK01Lora) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "atk01",
		Model:    "ATK01-LORA-01",
		Type:     transceivercom.LORA,
		Vendor:   "RHILEX-TECH",
		Mac:      "00:00:00:00:00:00",
		Firmware: "v0.0.0",
	}
}
func (tc *ATK01Lora) Status() transceivercom.TransceiverStatus {
	if tc.ProtocolSlaver == nil {
		return transceivercom.TransceiverStatus{
			Code:  transceivercom.TC_ERROR,
			Error: fmt.Errorf("Invalid Device"),
		}
	} else {
		return transceivercom.TransceiverStatus{
			Code:  transceivercom.TC_UP,
			Error: nil,
		}
	}

}
func (tc *ATK01Lora) Stop() {
	if tc.ProtocolSlaver != nil {
		tc.ProtocolSlaver.Stop()
	}
	glogger.GLogger.Info("ATK01-LORA Stopped")
}
