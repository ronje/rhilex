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
	"encoding/binary"
	"errors"
)

type GenericAppLayer struct {
	datalink *DataLinkLayer
}

func NewGenericAppLayerAppLayer(config TransporterConfig) *GenericAppLayer {
	return &GenericAppLayer{datalink: NewDataLinkLayer(config)}
}

func (app *GenericAppLayer) Request(appframe AppLayerFrame) (AppLayerFrame, error) {
	appBytes, errEncode := appframe.Encode()
	if errEncode != nil {
		return AppLayerFrame{}, errEncode
	}
	bytes, errHd := app.datalink.DataLinkerLayerHandle(appBytes)
	if errHd != nil {
		return AppLayerFrame{}, errHd
	}
	if len(bytes) < 4 {
		return AppLayerFrame{}, errors.New("data too short for header")
	}
	var responseHeader Header
	copy(responseHeader.Type[:], bytes[:2])
	copy(responseHeader.Length[:], bytes[2:4])
	payloadLength := int(binary.BigEndian.Uint16(responseHeader.Length[:]))
	if len(bytes) < 4+payloadLength {
		return AppLayerFrame{}, errors.New("data too short for payload")
	}
	payload := bytes[4 : 4+payloadLength]
	return AppLayerFrame{
		Header:  responseHeader,
		Payload: payload,
	}, nil
}

func (app *GenericAppLayer) Status() error {
	return app.datalink.Status()
}
func (app *GenericAppLayer) Close() error {
	return app.datalink.Close()
}
