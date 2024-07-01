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
	DataBuffer []byte
}

func NewATK01Lora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &ATK01Lora{
		R:          R,
		locker:     sync.Mutex{},
		DataBuffer: make([]byte, 1024),
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
func (tc *ATK01Lora) Start(Config transceivercom.TransceiverConfig) error {
	env := os.Getenv("LORASUPPORT")
	if env == "ATK01" {
		{
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
				buffer := [1024]byte{}
				var acc int          // 计数器而不是下标
				edgeSignal1 := false // 两个边沿
				edgeSignal2 := false // 两个边沿
				Bytes := [1]byte{}   // 读一个字节出来
				for {
					select {
					case <-typex.GCTX.Done():
						return
					default:
					}
					N, errR := io.Read(Bytes[:])
					if errR != nil {
						if !strings.Contains(errR.Error(), "timeout") {
							glogger.GLogger.Error(errR)
							return
						}
					}
					if N == 0 {
						continue
					}
					currentByte := Bytes[0]
					buffer[acc] = currentByte
					if acc > 0 {
						if currentByte == '\xEF' && buffer[acc-1] == '\xEE' {
							edgeSignal1 = true
						}
						if currentByte == '\n' && buffer[acc-1] == '\r' {
							// 注意：CRC校验不包含包头和包尾 [......)
							dataStartPos := 2
							crcPosL := acc - 2
							crcPosH := acc - 3
							crcByte := [2]byte{buffer[crcPosH], buffer[crcPosL]}
							crcCheckedValue := uint16(crcByte[0])<<8 | uint16(crcByte[1])
							dataBytes := buffer[dataStartPos:crcPosH]
							crcCalculatedValue := utils.CRC16(dataBytes)
							if crcCalculatedValue == crcCheckedValue {
								edgeSignal2 = true
							}
						}
					}
					// ++______--++_______--++_______--
					if edgeSignal1 && edgeSignal2 {
						glogger.GLogger.Debug(buffer[:acc+1])
						acc = 0
						edgeSignal1 = false
						edgeSignal2 = false
					}
					acc++
				}
			}(tc.serialPort)
		}
		glogger.GLogger.Info("ATK01-LORA-SX1278 Started")
	}
	return nil
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
		Name:     "ATK01",
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
