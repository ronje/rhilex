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

package mbus

import (
	"io"

	"github.com/hootrhino/rhilex/typex"
)

type MbusSerialTransporter struct {
	port typex.GenericRWC
}

func (dlt *MbusSerialTransporter) SendFrame(aduRequest []byte) (aduResponse []byte, err error) {
	if _, err = dlt.port.Write(aduRequest); err != nil {
		return nil, err
	}
	return dlt.ReadFrame(dlt.port)
}

/**
 * 读请求
 *
 */
func (dlt *MbusSerialTransporter) ReadFrame(rwc io.ReadWriteCloser) (aduResponse []byte, err error) {
	return aduResponse, nil
}
