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

# InternalStore
## 简介
InternalStore是未来新版本InterDb的替代。 是一个基于 GORM 的数据存储模块，用于管理数据库操作。它提供了一系列的 API 来实现常见的数据库操作，如创建、读取、更新、删除（CRUD），以及数据迁移和分页查询等功能。该模块被设计为包级单例模式，确保在整个应用程序中只有一个数据库连接实例，从而提高资源利用率和数据一致性。