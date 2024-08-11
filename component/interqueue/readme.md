<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

## 内部队列
主要用来数据缓冲用
- xqueue：老版本的消息队列，用了Go内置的Channel作为缓冲队列，已经触发到其极限了。
- yqueue：新版本的消息队列，使用list.List实现，动态扩容但是可能会消耗内存。

代码简单就不做赘述，稍微读一下即可看懂。