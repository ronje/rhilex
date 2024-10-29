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

package userprotocol

import (
	"time"

	"github.com/hootrhino/rhilex/typex"
)

type TransporterConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Port         typex.GenericRWC
}

type UserProtocolTransporter struct {
	port         typex.GenericRWC
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewUserProtocolTransporter(config TransporterConfig) *UserProtocolTransporter {
	return &UserProtocolTransporter{
		readTimeout:  config.ReadTimeout,
		writeTimeout: config.WriteTimeout,
		port:         config.Port,
	}
}

func (dlt *UserProtocolTransporter) SendFrame(aduRequest []byte) ([]byte, error) {
	deadline := time.Now().Add(dlt.writeTimeout)
	dlt.port.SetWriteDeadline(deadline)
	if _, err := dlt.port.Write(aduRequest); err != nil {
		return nil, err
	}
	defer dlt.port.SetWriteDeadline(time.Time{})
	return dlt.ReadFrame(dlt.port)
}

func (dlt *UserProtocolTransporter) ReadFrame(rwc typex.GenericRWC) ([]byte, error) {
	aduResponse := [128]byte{}
	deadline := time.Now().Add(dlt.readTimeout)
	dlt.port.SetReadDeadline(deadline)
	N, err := dlt.port.Read(aduResponse[:])
	if err != nil {
		return aduResponse[:N], err
	}
	defer dlt.port.SetReadDeadline(time.Time{})
	return aduResponse[:N], nil
}
