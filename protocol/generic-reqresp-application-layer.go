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
	errTxCount int32 // 错误包计数器
	errRxCount int32 // 错误包计数器
	datalink   *DataLayer
}

func NewGenericAppLayerAppLayer(config TransporterConfig) *GenericAppLayer {
	return &GenericAppLayer{errTxCount: 0, errRxCount: 0, datalink: NewDataLayer(config)}
}

func (app *GenericAppLayer) Request(appframe AppLayerFrame) (AppLayerFrame, error) {
	errWrite := app.Write(appframe)
	if errWrite != nil {
		app.errTxCount++
		return AppLayerFrame{}, errWrite
	}
	responseFrame, errRead := app.Read()
	if errRead != nil {
		app.errRxCount++
		return AppLayerFrame{}, errRead
	}
	return responseFrame, nil
}

func (app *GenericAppLayer) Write(appframe AppLayerFrame) error {
	appBytes, errEncode := appframe.Encode()
	if errEncode != nil {
		app.errTxCount++
		return errEncode
	}
	return app.datalink.Write(appBytes)
}

func (app *GenericAppLayer) Read() (AppLayerFrame, error) {
	bytes, errHd := app.datalink.Read()
	if errHd != nil {
		app.errRxCount++
		return AppLayerFrame{}, errHd
	}
	if len(bytes) < 4 {
		app.errRxCount++
		return AppLayerFrame{}, errors.New("data too short for header")
	}
	var responseHeader Header
	copy(responseHeader.Type[:], bytes[:2])
	copy(responseHeader.Length[:], bytes[2:4])
	payloadLength := int(binary.BigEndian.Uint16(responseHeader.Length[:])) - 1 /*crc byte*/
	if len(bytes) < 4+payloadLength {
		return AppLayerFrame{}, errors.New("data too short for payload")
	}
	payload := bytes[4 : 4+payloadLength]
	return AppLayerFrame{
		Header:  responseHeader,
		Payload: payload,
	}, nil
}

func (app *GenericAppLayer) GetErrTxCount() int32 {
	return app.errTxCount + app.datalink.errTxCount
}
func (app *GenericAppLayer) GetErrRxCount() int32 {
	return app.errRxCount + app.datalink.errRxCount
}
func (app *GenericAppLayer) Status() error {
	return app.datalink.Status()
}

func (app *GenericAppLayer) Close() error {
	return app.datalink.Close()
}
