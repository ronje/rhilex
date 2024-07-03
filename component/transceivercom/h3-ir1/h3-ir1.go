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

package h3ir1

/*
*
* 红外线接收模块
$ ir-keytable
Found /sys/class/rc/rc0/ (/dev/input/event1) with:

	Name: sunxi-ir
	Driver: sunxi-ir, table: rc-empty
	lirc device: /dev/lirc0
	Supported protocols: other lirc rc-5 rc-5-sz jvc sony nec sanyo mce_kbd rc-6 sharp xmp
	Enabled protocols: lirc nec
	bus: 25, vendor/product: 0001:0001, version: 0x0100
	Repeat delay = 500 ms, repeat period = 125 ms
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

const __IR_DEV = "/dev/input/event1"

type timeval struct {
	Second  int32 `json:"second"`
	USecond int32 `json:"uSecond"`
}
type irInputEvent struct {
	Time  timeval `json:"-"`
	Type  uint16  `json:"-"`
	Code  uint16  `json:"code"`
	Value int32   `json:"value"`
}

func (v irInputEvent) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

type H3IR1Config struct {
	InputHandle string `json:"inputHandle"` // 信号源
}
type H3IR1 struct {
	R          typex.Rhilex
	mainConfig H3IR1Config
	irFd       *os.File
}

func NewH3IR1(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &H3IR1{R: R, mainConfig: H3IR1Config{
		InputHandle: __IR_DEV,
	}}
}
func (ird *H3IR1) Start(transceivercom.TransceiverConfig) error {
	glogger.GLogger.Info("H3IR1 Started")

	fd, err := os.Open(ird.mainConfig.InputHandle)
	if err != nil {
		fd.Close()
		return err
	}
	ird.irFd = fd
	go func(ird *H3IR1) {
		defer func() {
			ird.irFd.Close()
		}()
		buf := make([]byte, 1024)
		for {
			select {
			case <-typex.GCTX.Done():
				return
			default:
				{
				}
			}
			n1, e := ird.irFd.Read(buf)
			if e != nil {
				glogger.GLogger.Error(e)
				continue
			}
			if n1 > 0 {
				event := irInputEvent{}
				// (*[24]byte)(unsafe.Pointer(&event))[:]
				// buffer := [24]byte{}
				_, err := ird.irFd.Read((*[24]byte)(unsafe.Pointer(&event))[:])
				if err != nil {
					glogger.GLogger.Error(err)
					continue
				}
				glogger.GLogger.Infof(event.String())
			}
			time.Sleep(125 * time.Millisecond)
		}
	}(ird)
	return nil
}
func (ird *H3IR1) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (ird *H3IR1) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "ir1",
		Model:    "ir1-nec",
		Type:     transceivercom.IR,
		Vendor:   "NEC-IR",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (ird *H3IR1) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_ERROR,
		Error: fmt.Errorf("NOT SUPPORT"),
	}
}
func (ird *H3IR1) Stop() {
	glogger.GLogger.Info("H3IR1 Stopped")
}
