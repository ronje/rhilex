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

package szy2062016

// 下行
const DIR0 byte = 0

// 不分包
const DIV byte = 0

// FCB
const FCB byte = 0x30 // 00110000

// 起始帧
const CTRL_CODE_FRAME_START byte = 0x68

// 结束帧
const CTRL_CODE_FRAME_END byte = 0x16

// 定义参数编码的常量
const (
	FCCommand          byte = 0x00 // 0 命令
	FCRainfall         byte = 0x01 // 1 雨量参数
	FCWaterLevel       byte = 0x02 // 2 水位参数
	FCFlowRate         byte = 0x03 // 3 流量(水量)参数
	FCFlowSpeed        byte = 0x04 // 4 流速参数
	FCGatePosition     byte = 0x05 // 5 闸位参数
	FCPower            byte = 0x06 // 6 功率参数
	FCAirPressure      byte = 0x07 // 7 气压参数
	FCWindSpeed        byte = 0x08 // 8 风速参数
	FCWaterTemperature byte = 0x09 // 9 水温参数
	FCWaterQuality     byte = 0x0A // 10 水质参数
	FCSoilMoisture     byte = 0x0B // 11 土壤含水率参数
	FCEvaporation      byte = 0x0C // 12 蒸发量参数
	FCAlarmStatus      byte = 0x0D // 13 报警或状态参数
	FCComprehensive    byte = 0x0E // 14 综合参数
	FCWaterPressure    byte = 0x0F // 15 水压参数
)

func CtrlCodeCommand() byte {
	return SetControlCode(FCCommand)
}
func CtrlCodeRainfall() byte {
	return SetControlCode(FCRainfall)
}
func CtrlCodeWaterLevel() byte {
	return SetControlCode(FCWaterLevel)
}
func CtrlCodeFlowRate() byte {
	return SetControlCode(FCFlowRate)
}
func CtrlCodeFlowSpeed() byte {
	return SetControlCode(FCFlowSpeed)
}
func CtrlCodeGatePosition() byte {
	return SetControlCode(FCGatePosition)
}
func CtrlCodePower() byte {
	return SetControlCode(FCPower)
}
func CtrlCodeAirPressure() byte {
	return SetControlCode(FCAirPressure)
}
func CtrlCodeWindSpeed() byte {
	return SetControlCode(FCWindSpeed)
}
func CtrlCodeWaterTemperature() byte {
	return SetControlCode(FCWaterTemperature)
}
func CtrlCodeWaterQuality() byte {
	return SetControlCode(FCWaterQuality)
}
func CtrlCodeSoilMoisture() byte {
	return SetControlCode(FCSoilMoisture)
}
func CtrlCodeEvaporation() byte {
	return SetControlCode(FCEvaporation)
}
func CtrlCodeAlarmStatus() byte {
	return SetControlCode(FCAlarmStatus)
}
func CtrlCodeComprehensive() byte {
	return SetControlCode(FCComprehensive)
}
func FrameDivide() bool {
	return DIV == 1
}
func FrameDirection() byte {
	return DIR0
}
func FrameFCB() byte {
	return FCB
}

// 获取控制码
func SetControlCode(FC byte) byte {
	return DIR0 | DIV | FCB | FC
}

// 解析
func ParseControlCode(originalByte byte) [4]byte {
	var maskDIR0 byte = 0x01 << 7          // 7
	var maskDIV byte = 0x01 << 6           // 6
	var maskFCB byte = 0x03 << 4           // 5 4
	var maskCode byte = 0x1F << 3          // 3 2 1 0
	dir0 := (originalByte & maskDIR0) >> 7 // 右移7位获取DIR0
	div := (originalByte & maskDIV) >> 6   // 右移6位获取DIV
	fcb := (originalByte & maskFCB) >> 4   // 右移4位获取FCB
	code := (originalByte & maskCode) >> 3 // 右移3位获取Code
	return [4]byte{dir0, div, fcb, code}
}
