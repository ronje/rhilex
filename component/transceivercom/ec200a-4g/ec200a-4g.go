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

package ec200a4g

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

type EC200ADtuConfig struct {
	Address string
}
type EC200ADtu struct {
	R          typex.Rhilex
	mainConfig EC200ADtuConfig
}

func NewEC200ADtu(R typex.Rhilex) transceivercom.TransceiverCommunicator {
	return &EC200ADtu{R: R, mainConfig: EC200ADtuConfig{
		Address: "/dev/ttyUSB1",
	}}
}
func (tc *EC200ADtu) Start(config transceivercom.TransceiverConfig) error {
	env := os.Getenv("4GSUPPORT")
	if env == "EC200A" {
		glogger.GLogger.Info("EC200A Init 4G")
		InitEC200A4G(config.Address)
		glogger.GLogger.Info("EC200A Init 4G Ok.")
	}
	glogger.GLogger.Info("EC200ADtu Started")
	return nil
}

type CSQInfo struct {
	Cops  string `json:"cops"`
	Csq   int    `json:"csq"`
	ICCID string `json:"iccid"`
}

func (tc *EC200ADtu) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	glogger.GLogger.Debug("EC200ADtu.Ctrl=", topic, string(args))
	if string(topic) == "mn4g.ec200a.info.csq" {
		CSQInfo1 := CSQInfo{
			Cops:  "CMCC",
			Csq:   15,
			ICCID: "00000000",
		}
		bytes, _ := json.Marshal(CSQInfo1)
		return bytes, nil
	}
	if string(topic) == "mn4g.ec200a.opt.restart" {
		return []byte("OK"), nil
	}
	if string(topic) == "mn4g.ec200a.cmd.send" {
		return []byte("OK"), nil
	}
	return []byte("OK"), nil
}
func (tc *EC200ADtu) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:     "EC200A-4G-DTU",
		Model:    "EC200A-CAT4",
		Type:     transceivercom.MN4G,
		Vendor:   "Quectel technology",
		Mac:      "00:00:00:00:00:00:00:00",
		Firmware: "0.0.0",
	}
}
func (tc *EC200ADtu) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_DOWN,
		Error: nil,
	}
}
func (tc *EC200ADtu) Stop() {
	glogger.GLogger.Info("EC200ADtu Stopped")
}
