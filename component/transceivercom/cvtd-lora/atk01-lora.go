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

package CvtdLora

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"

	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type CvtdLoraConfig struct {
	ComConfig transceivercom.TransceiverConfig
}
type CvtdLora struct {
	R          typex.Rhilex
	locker     sync.Mutex
	mainConfig CvtdLoraConfig
	serialPort serial.Port
	DataBuffer chan []byte
}

func NewCvtdLora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &CvtdLora{
		R:          R,
		locker:     sync.Mutex{},
		DataBuffer: make(chan []byte, 102400),
		mainConfig: CvtdLoraConfig{
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
func (tc *CvtdLora) Start(Config transceivercom.TransceiverConfig) error {
	env := os.Getenv("LORASUPPORT")
	if env != "CVTDLORA" {
		// return nil
	}
	tc.mainConfig = CvtdLoraConfig{
		ComConfig: Config,
	}
	glogger.GLogger.Info("Cvtd-LORA-01 Init")
	serialPort, err := tc.startSerialPort()
	if err != nil {
		return err
	}
	tc.serialPort = serialPort
	go tc.startLoopReceive(serialPort)
	go tc.startPushPacket(tc.DataBuffer)
	glogger.GLogger.Info("Cvtd-LORA-01 Started")
	return nil
}

/*
*
* 打开串口
*
 */
func (tc *CvtdLora) startSerialPort() (serial.Port, error) {
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
func (tc *CvtdLora) startPushPacket(Chan chan []byte) {
	for {
		select {
		case <-typex.GCTX.Done():
			return
		case buffer := <-Chan:
			Len := len(buffer)
			if Len > 2 {
				glogger.GLogger.Debug("CvtdLora Received Data:", buffer)
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
					Event:   "transceiver.up.data.cvtd",
					Ts:      uint64(time.Now().UnixMilli()),
					Summary: "transceiver.up.data",
					Info:    buffer,
				})
			}
		}
	}
}
func (tc *CvtdLora) startLoopReceive(io io.ReadWriteCloser) {
	MAX_BUFFER_SIZE := 1024 * 10 * 10
	buffer := make([]byte, MAX_BUFFER_SIZE)
	byteACC := 0                      // 计数器而不是下标
	edgeSignal1 := false              // 两个边沿
	edgeSignal2 := false              // 两个边沿
	expectPacket := make([]byte, 256) // 默认最大包长256字节

	for {
		select {
		case <-typex.GCTX.Done():
			return
		default:
		}
		N, errR := io.Read(buffer[byteACC:])
		// 读取异常，重启
		if errR != nil {
			if !strings.Contains(errR.Error(), "timeout") {
				glogger.GLogger.Error(errR)
			START:
				select {
				case <-typex.GCTX.Done():
					return
				default:
				}
				time.Sleep(5 * time.Second)
				serialPort, err1 := tc.startSerialPort()
				if err1 != nil {
					glogger.GLogger.Error("start Serial Port failed, try to restart,", err1)
					goto START
				} else {
					glogger.GLogger.Info("try to restart Serial port")
					tc.serialPort = serialPort
					go tc.startLoopReceive(tc.serialPort)
				}

				return
			}
		}
		if N == 0 {
			continue
		}
		byteACC += N
		if byteACC > 256 { // 单个包最大256字节
			if !edgeSignal1 || !edgeSignal2 {
				glogger.GLogger.Error("maximum data packet length(256) exceeded!")
				for i := 0; i < byteACC; i++ {
					buffer[i] = '\x00'
				}
				byteACC = 0
				edgeSignal1 = false
				edgeSignal2 = false
			}
			continue
		}
		expectPacketACC := 0
		expectPacketLength := 0
		for i, currentByte := range buffer[:byteACC] {
			expectPacketACC++
			if !edgeSignal1 {
				if expectPacketACC >= 2 {
					if currentByte == 0xEF && buffer[i-1] == 0xEE {
						edgeSignal1 = true
					}
				}
			}
			if !edgeSignal1 || expectPacketACC < 4 {
				continue
			}
			if edgeSignal1 {
				if currentByte == 0x0A && buffer[expectPacketACC-2] == 0x0D {
					expectPacketLength = copy(expectPacket, buffer[:expectPacketACC-1])
					edgeSignal2 = true
				}
			}
			if !edgeSignal1 || !edgeSignal2 {
				continue
			}
			if edgeSignal1 && edgeSignal2 {
				tc.DataBuffer <- expectPacket[2 : expectPacketLength-1]
			}
			if expectPacketACC < byteACC {
				if !edgeSignal1 || !edgeSignal2 {
					copy(buffer[0:], buffer[expectPacketACC-1:byteACC])
					byteACC = byteACC - expectPacketACC
				}
			} else {
				byteACC = 0
			}
			expectPacketLength = 0
			expectPacketACC = 0
			edgeSignal1 = false
			edgeSignal2 = false
		}
	}
}
func (tc *CvtdLora) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	if string(topic) == "lora.cvtd.cmd.send" {
		if tc.serialPort != nil {
			tc.serialPort.Write(args)
		}
		return []byte("OK"), nil
	}
	return []byte("OK"), nil
}
func (tc *CvtdLora) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "cvtd",
		Model:    "CVTD-LORA-01",
		Type:     transceivercom.LORA,
		Vendor:   "Beijing ChangWeiTongDa technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *CvtdLora) Status() transceivercom.TransceiverStatus {
	if tc.serialPort == nil {
		return transceivercom.TransceiverStatus{
			Code:  transceivercom.TC_ERROR,
			Error: fmt.Errorf("NOT SUPPORT"),
		}
	} else {
		return transceivercom.TransceiverStatus{
			Code:  transceivercom.TC_UP,
			Error: nil,
		}
	}

}
func (tc *CvtdLora) Stop() {
	if tc.serialPort != nil {
		tc.serialPort.Close()
	}
	glogger.GLogger.Info("Cvtd-LORA Stopped")
}
