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

package test

import (
	"fmt"
	"testing"

	"github.com/hootrhino/rhilex/device/dlt6452007"
)

// go test -timeout 30s -run ^TestCodec_DLT645_2007_Frame github.com/hootrhino/rhilex/test -v -count=1
// 68 45 92 66 23 00 10 68 11 04 33 34 34 35 25 16
// 68 45 92 66 23 00 10 68 91 06 33 34 34 35 66 55 62 16
// ============================= 33 34 34 35 66 55
// ============================= 00 01 01 02 33 22
func TestCodec_DLT645_2007_Frame(t *testing.T) {
	frame := dlt6452007.DLT645Frame0x11{
		Start:      dlt6452007.CTRL_CODE_FRAME_START,
		Address:    []byte{0x45, 0x92, 0x66, 0x23, 0x00, 0x10},
		CtrlCode:   dlt6452007.CTRL_CODE_READ_DATA,
		DataLength: 0x04,
		DataType:   [4]byte{0x33, 0x34, 0x34, 0x35},
		End:        dlt6452007.CTRL_CODE_FRAME_END,
	}
	packedFrame, err := frame.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("frame.Pack():")
	println()
	for _, v := range packedFrame {
		fmt.Printf(" %x", v)
	}
	println()
	responseFrame, err := dlt6452007.DecodeDLT645Frame0x11([]byte{
		0x68,
		0x45, 0x92, 0x66, 0x23, 0x00, 0x10,
		0x68,
		0x91,
		0x06, 0x33, 0x34, 0x34, 0x35, 0x66, 0x55, 0x62,
		0x16,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Unpacked frame: %+v\n", responseFrame.String())
	Data1, err1 := responseFrame.GetDataType()
	if err1 != nil {
		panic(err1)
	}
	t.Log(Data1)
	Data2, err2 := responseFrame.GetData()
	if err2 != nil {
		panic(err2)
	}
	t.Log(Data2)
}
