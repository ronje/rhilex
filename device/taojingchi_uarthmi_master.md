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

# 陶晶驰串口屏接入
主要对陶晶驰串口屏做了轻量支持，方便快捷接入到RHILEX。没有实现复杂逻辑，仅仅做了基本的串口指令包装。

## 示例
```lua
local err = tjchmi:Ctrl("Write", "t0.txt=\"Hello\"")
```