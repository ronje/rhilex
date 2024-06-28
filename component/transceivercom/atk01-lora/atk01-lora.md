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
# ATK-LORA-01
![atk01](image/atk01-lora/1719544633526.png)
ATK-LORA-01 一款体积小、微功率、低功耗、高性能远距离 LORA 无线串口模块。模块设计是采用高效的 ISM 频段射频 SX1278 扩频芯片，模块的工作频率在 410Mhz~441Mhz，以 1Mhz 频率为步进信道，共 32 个信道。可通过 AT 指令在线修改串口速率，发射功率，空中速率，工作模式等各种参数，并且支持固件升级功能。

## Topic
| Topic               | Args | Result(HTTP返回JSON里的data字段) |
| ------------------- | ---- | -------------------------------- |
| lora.atk01.cmd.send | ""   | String                           |
