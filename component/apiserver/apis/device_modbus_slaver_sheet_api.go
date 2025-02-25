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
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/typex"
)

func InitModbusSlaverRoute() {
	modbusApi := server.RouteGroup(server.ContextUrl("/modbus_slaver_sheet"))
	{
		modbusApi.GET(("/list"), server.AddRoute(ModbusSlaverSheetPageList))
	}
}

type ModbusSlaverRegister struct {
	UUID string `json:"uuid"`
	// 1: 离散输出Coils      Discrete Outputs
	// 2: 离散输入Coils      Discrete Inputs
	// 3: 保持寄存器         Holding Registers
	// 4: 输入寄存器         Input Registers
	Type    int `json:"type"`
	Address int `json:"address"`
	Value   any `json:"value"`
}
type ModbusSlaverRegisterVo struct {
	Coils             []ModbusSlaverRegister `json:"coils"`
	DiscreteRegisters []ModbusSlaverRegister `json:"discreteRegisters"`
	HoldingRegisters  []ModbusSlaverRegister `json:"holdingRegisters"`
	InputRegisters    []ModbusSlaverRegister `json:"inputRegisters"`
}
type RegisterEntityVo struct {
	UUID            string `json:"uuid"`
	AddressCoils    int    `json:"addressCoils"`
	ValueCoils      string `json:"valueCoils"`
	AddressDiscrete int    `json:"addressDiscrete"`
	ValueDiscrete   string `json:"valueDiscrete"`
	AddressHolding  int    `json:"addressHolding"`
	ValueHolding    string `json:"valueHolding"`
	AddressInput    int    `json:"addressInput"`
	ValueInput      string `json:"valueInput"`
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
	DiscreteRegisters := []ModbusSlaverRegister{}
	AllList := []RegisterEntityVo{}
	{
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
	}
	{
		for i := 0; i < 64; i++ {
			UUID := fmt.Sprintf("%s_DiscreteRegisters:%d", deviceUuid, i)
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
			DiscreteRegisters = append(DiscreteRegisters, Register)
		}
	}
	{
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
			HoldingRegisters = append(HoldingRegisters, Register)
		}
	}
	{
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
			InputRegisters = append(InputRegisters, Register)
		}
		for i := 0; i < 64; i++ {
			RegisterEntityVo := RegisterEntityVo{
				UUID:            fmt.Sprintf("uuid:%d", i),
				AddressCoils:    i,
				ValueCoils:      transValue(Coils[i].Value),
				AddressInput:    i,
				ValueInput:      transValue(InputRegisters[i].Value),
				AddressDiscrete: i,
				ValueDiscrete:   transValue(DiscreteRegisters[i].Value),
				AddressHolding:  i,
				ValueHolding:    transValue(HoldingRegisters[i].Value),
			}

			AllList = append(AllList, RegisterEntityVo)
		}
	}
	Result := service.WrapPageResult(*pager, AllList, 64)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}
func transValue(V any) string {
	switch T := V.(type) {
	case string:
		return T
	case int:
		return fmt.Sprintf("%d", T)
	case int64:
		return fmt.Sprintf("%d", T)
	case uint64:
		return fmt.Sprintf("%d", T)
	}
	return "0"
}
