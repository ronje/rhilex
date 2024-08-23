package rhilexlib

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/glogger"
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
func F1(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function2
*
 */
func F2(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function3
*
 */
func F3(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {

		return 1
	}
}

/*
*
* Modbus Function4
*
 */
func F4(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
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

func F5(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Values := l.ToString(5)
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		for _, v := range HexValues {
			if v > 1 {
				l.Push(lua.LString("Value Only Support '00' or '01'"))
				return 1
			}
		}
		Device := rx.GetDevice(devUUID)
		if Device == nil {
			l.Push(lua.LString("Device is not exists"))
			return 1
		}

		if Device.Type != typex.GENERIC_MODBUS_MASTER {
			l.Push(lua.LString("Only support GENERIC_MODBUS device"))
			return 1
		}
		if Device.Device.Status() != typex.DEV_UP {
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		args, _ := json.Marshal([]common.RegisterW{
			{
				Function: 5,
				SlaverId: byte(slaverId),
				Address:  uint16(Address),
				Values:   HexValues,
			},
		})
		_, err0 := Device.Device.OnWrite([]byte("F5"), args)
		if err0 != nil {
			l.Push(lua.LString(err0.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
*     local error = modbus:F6("uuid1", 0, 1, "0001020304")

*
 */
func F6(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Values := l.ToString(5) // 必须是单个字节: 000100010001
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device == nil {
			l.Push(lua.LString("Device is not exists"))
			return 1
		}

		if Device.Type != typex.GENERIC_MODBUS_MASTER {
			l.Push(lua.LString("Only support GENERIC_MODBUS device"))
			return 1
		}
		if Device.Device.Status() != typex.DEV_UP {
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		args, _ := json.Marshal(common.RegisterW{
			Function: 6,
			SlaverId: byte(slaverId),
			Address:  uint16(Address),
			Quantity: uint16(1), //2字节
			Values:   HexValues,
		})
		_, err0 := Device.Device.OnWrite([]byte("F6"), args)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			l.Push(lua.LString(err0.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
  - Modbus Function15
    local error = modbus:F15("uuid1", 0, 1, "0001020304")

*
*/
func F15(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Quantity := l.ToNumber(5) // 必须是单个字节: 000100010001
		Values := l.ToString(6)   // 必须是单个字节: 000100010001
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device == nil {
			l.Push(lua.LString("Device is not exists"))
			return 1
		}
		if Device.Type != typex.GENERIC_MODBUS_MASTER {
			l.Push(lua.LString("Only support GENERIC_MODBUS device"))
			return 1
		}
		if Device.Device.Status() != typex.DEV_UP {
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		args, _ := json.Marshal(common.RegisterW{
			Function: 15,
			SlaverId: byte(slaverId),
			Address:  uint16(Address),
			Quantity: uint16(Quantity),
			Values:   HexValues,
		})
		_, err0 := Device.Device.OnWrite([]byte("F15"), args)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			l.Push(lua.LString(err0.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* Modbus Function16
*    local error = modbus:F16("uuid1", 0, 1, "0001020304")
 */
func F16(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		slaverId := l.ToNumber(3)
		Address := l.ToNumber(4)
		Quantity := l.ToNumber(5) //
		Values := l.ToString(6)   //
		HexValues, err := hex.DecodeString(Values)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		Device := rx.GetDevice(devUUID)
		if Device == nil {
			l.Push(lua.LString("Device is not exists"))
			return 1
		}
		if Device.Type != typex.GENERIC_MODBUS_MASTER {
			l.Push(lua.LString("Only support GENERIC_MODBUS device"))
			return 1
		}
		if Device.Device.Status() != typex.DEV_UP {
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		args, _ := json.Marshal(common.RegisterW{
			Function: 16,
			SlaverId: byte(slaverId),
			Address:  uint16(Address),
			Quantity: uint16(Quantity),
			Values:   HexValues,
		})
		_, err0 := Device.Device.OnWrite([]byte("F16"), args)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			l.Push(lua.LString(err0.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* 解析Modbus报文
*
 */
func ParseModbusByte(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
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
