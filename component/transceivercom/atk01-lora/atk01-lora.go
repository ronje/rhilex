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
	DataBuffer chan []byte
}

func NewATK01Lora(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &ATK01Lora{
		R:          R,
		locker:     sync.Mutex{},
		DataBuffer: make(chan []byte, 1024),
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
	if env != "ATK01" {
		// return nil
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
		cursor := 0          // 用来标记数据当前读到哪里了
		edgeSignal1 := false // 两个边沿
		edgeSignal2 := false // 两个边沿
		dataStartPos := 0    // 0xEE ->
		dataEndPos := 0      // <- \n
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
			if byteACC > MAX_BUFFER_SIZE {
				glogger.GLogger.Error("exceeds the maximum buffer size")
				if !edgeSignal1 || !edgeSignal2 {
					byteACC = 0
					edgeSignal1 = false
					edgeSignal2 = false
				}
				continue
			}

			for i := cursor; i < byteACC; i++ {
				currentByte := buffer[i]
				cursor++
				if byteACC >= 2 && !edgeSignal1 {
					if currentByte == 0xEF && buffer[i-1] == 0xEE {
						edgeSignal1 = true
						dataStartPos = i
					}
				}
				if byteACC > 4 {
					if !edgeSignal1 {
						continue
					}
					if currentByte == 0x0A && buffer[i-1] == 0x0D {
						if edgeSignal1 {
							dataEndPos = i
							crcL := dataEndPos - 2
							crcH := crcL - 1
							crcByte := [2]byte{buffer[crcH], buffer[crcL]}
							crcCheckedValue := uint16(crcByte[0])<<8 | uint16(crcByte[1])
							currentPkt := buffer[dataStartPos+1 : dataEndPos-3]
							crcCalculatedValue := utils.CRC16(currentPkt)
							if crcCalculatedValue == crcCheckedValue {
								edgeSignal2 = true
							} else {
								glogger.GLogger.Errorf("CRC Check Error: (Checked=%d,Calculated=%d), data=%v",
									crcCheckedValue, crcCalculatedValue, currentPkt)
								// byteACC -= cursor // 出错后丢包
								// cursor -= cursor  // 出错后丢包
								edgeSignal1 = false
								edgeSignal2 = false
							}
						}
					}
					if edgeSignal1 && edgeSignal2 {
						tc.DataBuffer <- buffer[dataStartPos+1 : dataEndPos-3]
						// glogger.GLogger.Debug(buffer[dataStartPos+1 : dataEndPos-3])
						if cursor >= byteACC {
							byteACC = 0
							cursor = 0
						} else {
							offset := byteACC - cursor
							copy(buffer[0:], buffer[offset:])
							byteACC += offset
						}
						edgeSignal1 = false
						edgeSignal2 = false
					}
				}
			}
		}
	}(tc.serialPort)
	go func(Chan chan []byte) {
		for {
			select {
			case <-typex.GCTX.Done():
				return
			case Data := <-Chan:
				glogger.GLogger.Debug(Data)
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
