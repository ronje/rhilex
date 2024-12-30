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

package protocol

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// go test -timeout 30s  -run ^TestGenericProtocolMaster$ github.com/hootrhino/rhilex/protocol -v -count=1
func TestGenericProtocolMaster(t *testing.T) {
	Logger := logrus.StandardLogger()
	Logger.SetLevel(logrus.DebugLevel)
	config := ExchangeConfig{
		Port:         NewGenericReadWriteCloser(),
		ReadTimeout:  5000,
		WriteTimeout: 5000,
		PacketEdger: PacketEdger{
			Head: [2]byte{0xAB, 0xAB},
			Tail: [2]byte{0xBA, 0xBA},
		},
		Logger: Logger,
	}
	TransportMaster := NewGenericProtocolMaster(config)
	Request := NewApplicationFrame([]byte{0xFF, 0xFF, 0xFF, 0xFF})
	t.Log("Request:", Request.ToString())
	Response, err := TransportMaster.Request(Request)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Response:", Response.ToString())
	}
}
