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

package iec104

import "fmt"

// APDU 104数据包
type APDU struct {
	APCI     *APCI
	ASDU     *ASDU
	Len      int
	ASDULen  int
	CtrType  byte
	CtrFrame any
	Signals  []*Signal
}

// parseAPDU 解析APDU
func (apdu *APDU) parseAPDU(input []byte) error {
	if input == nil || len(input) < 4 {
		return fmt.Errorf("APDU报文[%X]非法", input)
	}
	apci := &APCI{
		ApduLen: len(input),
		Ctr1:    input[0],
		Ctr2:    input[1],
		Ctr3:    input[2],
		Ctr4:    input[3],
	}
	fType, ctrFrame, err := apci.ParseCtr()
	if err != nil {
		return fmt.Errorf("APDU报文[%X]解析控制域异常: %v", input, err)
	}
	asdu := new(ASDU)
	var asduLen int
	signals := make([]*Signal, 0)
	if len(input[4:]) < 1 {
		asduLen = 0
	} else {
		signals, err = asdu.ParseASDU(input[4:])
		if err != nil {
			return fmt.Errorf("APDU报文[%X]解析ASDU域[%X]异常: %v", input, input[4:], err)
		}
		asduLen = len(input[6:])
	}
	apdu.APCI = apci
	apdu.ASDU = asdu
	apdu.Len = apci.ApduLen
	apdu.ASDULen = asduLen
	apdu.CtrType = fType
	apdu.CtrFrame = ctrFrame
	apdu.Signals = signals
	return nil
}
