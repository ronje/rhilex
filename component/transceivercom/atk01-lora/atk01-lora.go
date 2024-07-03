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
	env := os.Getenv("LORASUPPORT")
	if env != "ATK01" {
		return nil
	}
	glogger.GLogger.Info("ATK01-LORA-SX1278 Init")
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
	tc.serialPort = serialPort
	go func(io io.ReadWriteCloser) {
		MAX_BUFFER_SIZE := 1024 * 10 * 10
		buffer := make([]byte, MAX_BUFFER_SIZE)
		byteACC := 0         // 计数器而不是下标
		edgeSignal1 := false // 两个边沿
		edgeSignal2 := false // 两个边沿
		expectPacket := make([]byte, 256)

		for {
			select {
			case <-typex.GCTX.Done():
				return
			default:
			}
			N, errR := io.Read(buffer[byteACC:])
			if errR != nil {
				if !strings.Contains(errR.Error(), "timeout") {
					glogger.GLogger.Error(errR)
					return
				}
			}
			if N == 0 {
				continue
			}
			byteACC += N
			if byteACC > 256 { // 单个包最大256字节
				if !edgeSignal1 || !edgeSignal2 {
					glogger.GLogger.Error("maximum data packet length(256) exceeded, will flush all data!")
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
	}(tc.serialPort)
	go func(Chan chan []byte) {
		for {
			select {
			case <-typex.GCTX.Done():
				return
			case buffer := <-Chan:
				glogger.GLogger.Debug("Received:", buffer)
				Len := len(buffer)
				if Len > 2 {
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
	}(tc.DataBuffer)
	glogger.GLogger.Info("ATK01-LORA-SX1278 Started")
	return nil
}

/*
*
* 递归解析
*
 */
func ParsePacket(binary []byte, offset int) {
	if len((binary)) == 0 {
		return
	}
	N := 0
	ParsePacket(binary[offset+N:], offset+N)
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
		Model:    "ATK01-LORA-SX1278",
		Type:     transceivercom.LORA,
		Vendor:   "GUANGZHOU-ZHENGDIAN-YUANZI technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *ATK01Lora) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_ERROR,
		Error: fmt.Errorf("NOT SUPPORT"),
	}
}
func (tc *ATK01Lora) Stop() {
	if tc.serialPort != nil {
		tc.serialPort.Close()
	}
	glogger.GLogger.Info("EC200ADtu Stopped")
}
func ShiftBytes(N int, Bytes *[]byte) {
	if Bytes == nil || len(*Bytes) == 0 || N == 0 {
		return
	}
	N = N % len(*Bytes)
	bytes := *Bytes
	front := bytes[N:]
	back := bytes[:N]
	*Bytes = append(front, back...)
}
