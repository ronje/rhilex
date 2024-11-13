<!--
 Copyright (C) 2024 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# GPS 信号接收网关
GPS模块发出的信号通常是以NMEA（National Marine Electronics Association）协议格式传输的，这是一种通用的数据传输协议。以下是一个典型的NMEA 0183协议格式的GPS数据示例：
```
$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47
```
这个数据字符串由多个字段组成，每个字段都有特定的含义：
- `$GPGGA`: 固定句头，表示Global Positioning System Fix Data（全球定位系统固定数据）
- `123519`: UTC时间，格式为hhmmss（小时、分钟、秒）
- `4807.038`: 纬度，格式为ddmm.mmm（度、分、分的小数部分）
- `N`: 纬度半球，N表示北半球，S表示南半球
- `01131.000`: 经度，格式为dddmm.mmm（度、分、分的小数部分）
- `E`: 经度半球，E表示东经，W表示西经
- `1`: 定位质量指示，0=未定位，1=非差分定位，2=差分定位
- `08`: 使用中的卫星数量（00-12）
- `0.9`: 水平精确度（HDOP）
- `545.4`: 海拔高度，单位为米
- `M`: 海拔高度单位，M表示米
- `46.9`: 大地水准面高度，单位为米
- `M`: 大地水准面高度单位，M表示米
- `,,`: 空字段
- `*47`: 校验和，由星号和两个十六进制数字组成，用于验证数据完整性
请注意，这只是GPS模块可能发出的众多NMEA句子之一。其他常见的NMEA句子包括 `$GPRMC`（推荐最小定位信息）、`$GPGLL`（地理定位信息）等。每个句子都以美元符号（`$`）开始，以回车换行符（`\r\n`）结束。
