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

type GenericAppLayer struct {
	errTxCount int32 // 错误包计数器
	errRxCount int32 // 错误包计数器
	transport  *TransportLayer
}

func NewGenericAppLayerAppLayer(config TransporterConfig) *GenericAppLayer {
	return &GenericAppLayer{errTxCount: 0, errRxCount: 0, transport: NewTransportLayer(config)}
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
	return app.transport.Write(appBytes)
}

func (app *GenericAppLayer) Read() (AppLayerFrame, error) {
	bytes, errHd := app.transport.Read()
	if errHd != nil {
		app.errRxCount++
		return AppLayerFrame{}, errHd
	}
	Frame, errDecode := DecodeAppLayerFrame(bytes)
	if errDecode != nil {
		app.errRxCount++
		return AppLayerFrame{}, errDecode
	}
	return Frame, nil
}

func (app *GenericAppLayer) Status() error {
	return app.transport.Status()
}

func (app *GenericAppLayer) Close() error {
	return app.transport.Close()
}
