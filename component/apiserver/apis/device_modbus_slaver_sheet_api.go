// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/typex"
)

type ModbusSlaverRegister struct {
	UUID string `json:"uuid"`
	// 1: 离散输出Coils      Discrete Outputs
	// 2: 离散输入Coils      Discrete Inputs
	// 3: 保持寄存器         Holding Registers
	// 4: 输入寄存器         Input Registers
	Type    int         `json:"type"`
	Address int         `json:"address"`
	Value   interface{} `json:"value"`
}
type ModbusSlaverRegisterVo struct {
	Coils            []ModbusSlaverRegister `json:"coils"`
	HoldingRegisters []ModbusSlaverRegister `json:"holdingRegisters"`
	InputRegisters   []ModbusSlaverRegister `json:"inputRegisters"`
}

func ModbusSlaverSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	Slot := intercache.GetSlot(deviceUuid)
	if Slot == nil {
		c.JSON(common.HTTP_OK, common.Error("Cache Slot Not Exists"))
		return
	}
	// 1: 线圈寄存器      Coils Registers
	// 2: 离散寄存器      Discrete Registers
	// 3: 保持寄存器      Holding Registers
	// 4: 输入寄存器      Input Registers
	Coils := []ModbusSlaverRegister{}
	HoldingRegisters := []ModbusSlaverRegister{}
	InputRegisters := []ModbusSlaverRegister{}
	for i := 0; i < 64; i++ {
		UUID := fmt.Sprintf("%s_Coils:%d", deviceUuid, i)
		Register := ModbusSlaverRegister{
			UUID:    UUID,
			Type:    1,
			Address: i,
			Value:   0,
		}
		Value, ok := Slot[UUID]
		if ok {
			Register.Value = Value.Value
		}
		Coils = append(Coils, Register)
	}
	for i := 0; i < 64; i++ {
		UUID := fmt.Sprintf("%s_HoldingRegisters:%d", deviceUuid, i)
		Register := ModbusSlaverRegister{
			UUID:    UUID,
			Type:    1,
			Address: i,
			Value:   0,
		}
		Value, ok := Slot[UUID]
		if ok {
			Register.Value = Value.Value
		}
		HoldingRegisters = append(Coils, Register)
	}
	for i := 0; i < 64; i++ {
		UUID := fmt.Sprintf("%s_InputRegisters:%d", deviceUuid, i)
		Register := ModbusSlaverRegister{
			UUID:    UUID,
			Type:    1,
			Address: i,
			Value:   0,
		}
		Value, ok := Slot[UUID]
		if ok {
			Register.Value = Value.Value
		}
		InputRegisters = append(Coils, Register)
	}
	Result := service.WrapPageResult(*pager, ModbusSlaverRegisterVo{
		Coils:            Coils,
		HoldingRegisters: HoldingRegisters,
		InputRegisters:   InputRegisters,
	}, 64)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}
