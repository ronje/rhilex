# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
import crcmod
import binascii

def packet(data):
    # 定义协议的起始和结束标记
    start_mark = b'\xEE\xEF'
    end_mark = b'\x0D\x0A'

    # 使用Modbus格式的CRC多项式计算CRC16校验码
    crc16_func = crcmod.mkCrcFun(0x18005, rev=True, initCrc=0xFFFF, xorOut=0x0000)
    crc16 = crc16_func(data)
    crc16_bytes = crc16.to_bytes(2, 'big')  # 将CRC16转换为2个字节的byte数组

    # 组合数据
    packet_data = start_mark + data + crc16_bytes + end_mark

    # 将结果转换为十六进制字符串
    hex_output = binascii.hexlify(data=packet_data,sep=",").decode('utf-8')

    return hex_output

# 示例使用
data1 = b'\x00'
data2 = b'\x01\x02'
data3 = b'\x01\x02\x03'
data4 = b'\x02\x03\x04\x05'
data5 = b'\x02\x03\x04'
data6 = b'\x02\x03'
data7 = b'\x09'
print(f"Generate Packet: [{packet(data1)}]")
print(f"Generate Packet: [{packet(data2)}]")
print(f"Generate Packet: [{packet(data3)}]")
print(f"Generate Packet: [{packet(data4)}]")
print(f"Generate Packet: [{packet(data5)}]")
print(f"Generate Packet: [{packet(data6)}]")
print(f"Generate Packet: [{packet(data7)}]")

