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
	"fmt"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"

	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/protocol"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type ATK01LoraConfig struct {
	ComConfig transceivercom.TransceiverConfig
}
type ATK01Lora struct {
	R          typex.Rhilex
	locker     sync.Mutex
	mainConfig ATK01LoraConfig
	serialPort serial.Port
	DataBuffer chan []byte
}

func NewATK01Lora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &ATK01Lora{
		R:          R,
		locker:     sync.Mutex{},
		DataBuffer: make(chan []byte, 102400),
		mainConfig: ATK01LoraConfig{
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
	serialPort, err := tc.startSerialPort()
	if err != nil {
		return err
	}
	tc.serialPort = serialPort
	go protocol.StartEdgeSymReceive(typex.GCTX, tc.DataBuffer, tc.serialPort)
	go tc.startProcessPacket(tc.DataBuffer)
	glogger.GLogger.Info("ATK01-LORA-01 Started")
	return nil
}

/*
*
* 打开串口
*
 */
func (tc *ATK01Lora) startSerialPort() (serial.Port, error) {
	config := serial.Config{
		Address:  tc.mainConfig.ComConfig.Address,
		BaudRate: tc.mainConfig.ComConfig.BaudRate,
		DataBits: tc.mainConfig.ComConfig.DataBits,
		Parity:   tc.mainConfig.ComConfig.Parity,
		StopBits: tc.mainConfig.ComConfig.StopBits,
		Timeout:  time.Duration(tc.mainConfig.ComConfig.IOTimeout) * time.Millisecond,
		RS485:    serial.RS485Config{},
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		return nil, err
	}
	return serialPort, nil
}

/*
*
* 向RHILEX推数据
*
 */
func (tc *ATK01Lora) startProcessPacket(Chan chan []byte) {
	for {
		select {
		case <-typex.GCTX.Done():
			return
		case buffer := <-Chan:
			Len := len(buffer)
			if Len > 2 {
				glogger.GLogger.Debug("ATK01Lora Received Data:", buffer)
				crcByte := [2]byte{buffer[Len-2], buffer[Len-1]}
				crcCheckedValue := uint16(crcByte[0])<<8 | uint16(crcByte[1])
				crcCalculatedValue := utils.CRC16(buffer[:Len-2])
				if crcCalculatedValue != crcCheckedValue {
					glogger.GLogger.Errorf("CRC Check Error: (Checked=%d,Calculated=%d), data=%v",
						crcCheckedValue, crcCalculatedValue, buffer)
					continue
				}
				internotify.Push(internotify.BaseEvent{
					Type:    "transceiver.up.data",
					Event:   "transceiver.up.data.atk01",
					Ts:      uint64(time.Now().UnixMilli()),
					Summary: "transceiver.up.data",
					Info:    buffer,
				})
			}
		}
	}
}

func (tc *ATK01Lora) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	if string(topic) == "lora.atk01.cmd.send" {
		if tc.serialPort != nil {
			tc.serialPort.Write(args)
		}
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
	if tc.serialPort == nil {
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
	if tc.serialPort != nil {
		tc.serialPort.Close()
	}
	glogger.GLogger.Info("ATK01-LORA Stopped")
}
