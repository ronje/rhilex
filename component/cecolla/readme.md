<!--
 Copyright (C) 2025 wwhai

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
# CecollaResourceManager
历史原因留下了不合理的设计：
```go
// 云边协同
CheckCecollaType(Type CecollaType) error
GetCecolla(string) *Cecolla
SaveCecolla(*Cecolla)
AllCecollas() []*Cecolla
RestartCecolla(uuid string) error
RemoveCecolla(uuid string)
LoadCecollaWithCtx(cecolla *Cecolla, ctx context.Context, cancelCTX context.CancelFunc) error
```
0.8重构CecollaResourceManager，从RHILEX里面拿出来作为组件.最终和多媒体管理器一致。
