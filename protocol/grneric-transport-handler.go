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

type GenericProtocolHandler struct {
	appLayer *GenericAppLayer
}

func NewGenericProtocolHandler(config TransporterConfig) *GenericProtocolHandler {
	return &GenericProtocolHandler{appLayer: NewGenericAppLayerAppLayer(config)}
}
func (handler *GenericProtocolHandler) Request(appframe AppLayerFrame) (AppLayerFrame, error) {
	return handler.appLayer.Request(appframe)
}
func (handler *GenericProtocolHandler) Status() error {
	return handler.appLayer.Status()
}
func (handler *GenericProtocolHandler) Close() error {
	return handler.appLayer.Close()
}
