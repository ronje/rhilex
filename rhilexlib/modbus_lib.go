package rhilexlib

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

//  --------------------------------------------
// |Function | Register Type
//  --------------------------------------------
// |	1	 | Read Coil
// |	2	 | Read Discrete Input
// |	3	 | Read Holding Registers
// |	4	 | Read Input Registers
// |	5	 | Write Single Coil
// |	6	 | Write Single Holding Register
// |	15	 | Write Multiple Coils
// |	16	 | Write Multiple Holding Registers
//  --------------------------------------------
/*
*
* Modbus Function1
*
 */
func F1(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function2
*
 */
func F2(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function3
*
 */
func F3(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function4
*
 */
func F4(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
  - Modbus Function5
    local error = modbus:F5("uuid1", 0, 1, "0001020304")

*
*/

func F5(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
*     local error = modbus:F6("uuid1", 0, 1, "0001020304")

*
 */
func F6(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
  - Modbus Function15
    local error = modbus:F15("uuid1", 0, 1, "0001020304")

*
*/
func F15(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function16
*    local error = modbus:F16("uuid1", 0, 1, "0001020304")
 */
func F16(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* 解析Modbus报文
*
 */
func ParseModbusByte(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		hexS := l.ToString(2)
		b, err0 := hex.DecodeString(hexS)
		if err0 != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err0.Error()))
			return 2
		}
		modbus, ok := parseModbusByte(b)
		if ok {
			T := lua.LTable{}
			T.RawSet(lua.LString(modbus.Address), lua.LNumber(modbus.Address))
			T.RawSet(lua.LString(modbus.Function), lua.LNumber(modbus.Function))
			Data := lua.LTable{}
			for _, v := range modbus.Data {
				Data.Append(lua.LNumber(v))
			}
			T.RawSet(lua.LString(modbus.Data), &Data)
			l.Push(&T)
			l.Push(lua.LNil)
			return 2
		}
		l.Push(lua.LNil)
		l.Push(lua.LString("parse modbus error"))
		return 2
	}
}

type ModbusData struct {
	Address  uint8
	Function uint8
	Data     []byte
}

func parseModbusByte(b []byte) (ModbusData, bool) {
	if len(b) < 4 {
		return ModbusData{}, false
	}
	address := b[0]
	function := b[1]
	data := b[2 : len(b)-2]
	receivedCRC := binary.BigEndian.Uint16(b[len(b)-2:])
	calculatedCRC := calculateCRC(b[:len(b)-2])
	if receivedCRC != calculatedCRC {
		return ModbusData{}, false
	}
	modbusData := ModbusData{
		Address:  address,
		Function: function,
		Data:     data,
	}
	return modbusData, true
}
func calculateCRC(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 > 0 {
				crc >>= 1
				crc ^= 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// POST -> temp , 0x0001
type CtrlCmd struct {
	Tag   string `json:"tag"`   // 点位表的Tag
	Value string `json:"value"` // 写的值
}

func (O CtrlCmd) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

/**
 * 设备控制
 *
 */
func WriteToModbusSheetRegisterWithTag(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		args := stateStack.ToString(3)
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				ctrlCmd := CtrlCmd{}
				if errUnmarshal := json.Unmarshal([]byte(args), &ctrlCmd); errUnmarshal != nil {
					stateStack.Push(lua.LString(errUnmarshal.Error()))
					return 1
				}
				_, err := Device.Device.OnCtrl([]byte("WriteToModbusSheetRegisterWithTag"), []byte(ctrlCmd.String()))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}
