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

package transceiver

import (
	"context"
	"fmt"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/component/eventbus"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/protocol"
	"github.com/hootrhino/rhilex/typex"
)

type DefalutTransceiver struct {
	rhilex         typex.Rhilex
	locker         sync.Mutex
	mainConfig     TransceiverConfig
	Event          string
	EventType      string
	ProtocolSlaver *protocol.GenericProtocolSlaver
}

func NewTransceiver(rhilex typex.Rhilex) Transceiver {
	return &DefalutTransceiver{
		rhilex: rhilex,
		locker: sync.Mutex{},
		mainConfig: TransceiverConfig{
			Name:              "default_transceiver",
			Address:           "COM1",
			BaudRate:          9600,
			DataBits:          8,
			Parity:            "N",
			StopBits:          1,
			IOTimeout:         0,   // IOTimeout * time.Millisecond
			ATTimeout:         200, // ATRwTimeout * time.Millisecond
			TransportProtocol: 1,
		},
	}
}

/*
*
* Start
*
 */
func (tc *DefalutTransceiver) Start(Config TransceiverConfig) error {
	tc.mainConfig = Config
	tc.EventType = "transceiver.up.data"
	tc.Event = fmt.Sprintf("transceiver.up.data.%s", tc.mainConfig.Name)
	glogger.GLogger.Info("Transceiver Init:", tc.mainConfig.Name)
	config := serial.Config{
		Address:  tc.mainConfig.Address,
		BaudRate: tc.mainConfig.BaudRate,
		DataBits: tc.mainConfig.DataBits,
		Parity:   tc.mainConfig.Parity,
		StopBits: tc.mainConfig.StopBits,
		Timeout:  time.Duration(tc.mainConfig.IOTimeout) * time.Millisecond,
	}
	serialPort, errOpen := serial.Open(&config)
	if errOpen != nil {
		glogger.GLogger.Error("serial port start failed err:", errOpen, ", config:", config)

		return errOpen
	}
	TransporterConfig := protocol.ExchangeConfig{
		Port:         serialPort,
		ReadTimeout:  time.Duration(tc.mainConfig.IOTimeout * int64(time.Millisecond)),
		WriteTimeout: time.Duration(tc.mainConfig.IOTimeout * int64(time.Millisecond)),
		Logger:       glogger.Logrus,
	}
	ctx, cancel := context.WithCancel(context.Background())
	tc.ProtocolSlaver = protocol.NewGenericProtocolSlaver(ctx, cancel, TransporterConfig)
	go tc.ProtocolSlaver.StartLoop(func(AppLayerFrame protocol.AppLayerFrame, errRead error) {
		if errRead != nil {
			glogger.GLogger.Error(errRead)
			return
		}
		glogger.GLogger.Debug("Transceiver.ProtocolSlaver.Receive:", AppLayerFrame.ToString())
		buffer, _ := AppLayerFrame.Encode()
		lineS := "event.transceiver.data." + tc.mainConfig.Address
		eventbus.Publish(lineS, eventbus.EventMessage{
			Topic:   lineS,
			From:    "transceiver",
			Type:    "HARDWARE",
			Event:   lineS,
			Ts:      uint64(time.Now().UnixMilli()),
			Payload: buffer,
		})
	})
	glogger.GLogger.Info("Transceiver Started:", tc.mainConfig.Name)
	return nil
}

func (tc *DefalutTransceiver) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	return []byte("OK"), nil
}
func (tc *DefalutTransceiver) Info() CommunicatorInfo {
	return CommunicatorInfo{
		Name:     tc.mainConfig.Name,
		Model:    "RHILEX TRANSCEIVERCOM",
		Type:     LORA,
		Vendor:   "RHILEX-TECH",
		Mac:      "00:00:00:00:00:00",
		Firmware: "v0.0.1",
	}
}
func (tc *DefalutTransceiver) Status() TransceiverStatus {
	if tc.ProtocolSlaver == nil {
		return TransceiverStatus{
			Code:  TC_ERROR,
			Error: fmt.Errorf("Invalid Device"),
		}
	} else {
		return TransceiverStatus{
			Code:  TC_UP,
			Error: nil,
		}
	}

}
func (tc *DefalutTransceiver) Stop() {
	if tc.ProtocolSlaver != nil {
		tc.ProtocolSlaver.Stop()
	}
	glogger.GLogger.Info("Transceiver Stopped:", tc.mainConfig.Name)

}
