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

package TemplateCom

import (
	"time"

	"github.com/hootrhino/rhilex/typex"
)

type TemplateComConfig struct {
}
type Mx01BLE struct {
	R          typex.Rhilex
	mainConfig TemplateComConfig
}

func NewMx01BLE(R typex.Rhilex) transceivercom.transceivercommunicator {
	return &Mx01BLE{R: R, mainConfig: TemplateComConfig{}}
}
func (tc *TemplateCom) Start(map[string]any) error {
	return nil
}
func (tc *TemplateCom) Ctrl(cmd []byte, timeout time.Duration) ([]byte, error) {
	return nil, nil
}
func (tc *TemplateCom) Info() transceivercom.CommunicatorInfo {
	return transceivercom.CommunicatorInfo{
		Name:   "BLE-Module",
		Model:  "01",
		Type:   transceivercom.BLE,
		Vendor: "COMPANY A",
	}
}
func (tc *TemplateCom) Status() transceivercom.TransceiverStatus {
	return transceivercom.TransceiverStatus{
		Code:  transceivercom.TC_UP,
		Error: nil,
	}
}
func (tc *TemplateCom) Stop() {

}
