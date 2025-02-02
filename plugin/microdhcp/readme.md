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

# microdhcp 配置说明

`microdhcp` 是一个轻量级的 DHCP 服务插件，允许你为局域网内的设备动态分配 IP 地址，并且支持静态 IP 映射。本文档将详细介绍如何配置和使用 `microdhcp` 插件。

## 插件功能

- **动态 IP 地址分配**：为网络中的设备分配 IP 地址。
- **静态 IP 映射**：为特定的设备根据 MAC 地址指定固定的 IP 地址。
- **租期设置**：配置 IP 地址的租期，决定设备在不重新请求 IP 地址的情况下能够使用该地址的最大时间。
- **网关和 DNS 配置**：为 DHCP 客户端提供网关和 DNS 服务器地址。

## 配置说明

下面是 `microdhcp` 插件的详细配置项及其说明：

```ini
[plugin.microdhcp]
# 启用 microdhcp 插件
enable = true

# microdhcp 监听的 IP 地址。设置为 0.0.0.0 表示监听所有可用的网络接口
listen_ip = 0.0.0.0

# microdhcp 监听的端口号。默认使用端口 67
listen_port = 67

# 静态地址映射。为指定的设备根据 MAC 地址分配固定的 IP 地址。
# 格式：MAC=IP，多个映射用逗号分隔。例如：
# 00:11:22:33:44:55=192.168.1.100, 00:11:22:33:44:56=192.168.1.101
static_address_mapping = "00:11:22:33:44:55=192.168.1.100, 00:11:22:33:44:56=192.168.1.101"

# 可选配置：DHCP 动态分配的 IP 地址范围。
# 在此范围内，设备将从 DHCP 服务中动态获取 IP 地址。
# 例如：动态分配的 IP 范围从 192.168.1.200 到 192.168.1.250
dhcp_range_start = 192.168.1.200
dhcp_range_end = 192.168.1.250

# 可选配置：DHCP 租期时间，单位为秒。
# 例如：86400 秒即 24 小时，表示客户端每 24 小时需要重新请求 IP 地址。
dhcp_lease_time = 86400  # 24 小时

# 可选配置：DHCP 服务提供的 DNS 服务器地址。
# 例如：设置为 Google 的公共 DNS 服务器 8.8.8.8 和 8.8.4.4
dns_servers = "8.8.8.8, 8.8.4.4"

# 可选配置：网关地址，DHCP 客户端将使用此地址作为默认网关。
# 例如：设置网关地址为 192.168.1.1
gateway_ip = "192.168.1.1"
```

## 配置项详解

### `enable`
- **描述**：启用或禁用 `microdhcp` 插件。
- **类型**：布尔值（`true` 或 `false`）。
- **默认值**：`true`。
- **说明**：如果设置为 `false`，插件将不会运行。

### `listen_ip`
- **描述**：配置 `microdhcp` 监听的 IP 地址。
- **类型**：字符串（IP 地址）。
- **默认值**：`0.0.0.0`（表示监听所有可用的网络接口）。
- **说明**：如果你希望 `microdhcp` 仅监听特定的网络接口，可以在此处指定具体的 IP 地址。

### `listen_port`
- **描述**：配置 `microdhcp` 监听的端口号。
- **类型**：整数。
- **默认值**：`67`（DHCP 协议默认使用端口 67）。
- **说明**：通常无需更改此设置，除非你有特殊的端口需求。

### `static_address_mapping`
- **描述**：配置静态 IP 映射，为特定设备分配固定 IP 地址。
- **类型**：字符串（例如：`00:11:22:33:44:55=192.168.1.100`）。
- **默认值**：空字符串（表示没有静态地址映射）。
- **说明**：此项可以为空。如果配置了静态映射，`microdhcp` 将为指定的 MAC 地址分配固定的 IP 地址。

### `dhcp_range_start` 和 `dhcp_range_end`
- **描述**：配置 DHCP 服务分配 IP 地址的范围。
- **类型**：字符串（IP 地址）。
- **默认值**：无。
- **说明**：为动态分配的设备指定 IP 地址范围。例如，设置为 `192.168.1.200` 到 `192.168.1.250`，则在此范围内的设备将通过 DHCP 获得 IP 地址。

### `dhcp_lease_time`
- **描述**：配置 DHCP 租期时间，单位为秒。
- **类型**：整数。
- **默认值**：`86400`（即 24 小时）。
- **说明**：租期时间决定了设备使用其 IP 地址的最大时间。到期后，设备需要重新向 DHCP 服务器请求 IP 地址。

### `dns_servers`
- **描述**：配置 DHCP 客户端使用的 DNS 服务器地址。
- **类型**：字符串（多个 DNS 服务器地址，用逗号分隔）。
- **默认值**：无。
- **说明**：为客户端提供 DNS 服务器地址，客户端将使用这些 DNS 服务器进行域名解析。

### `gateway_ip`
- **描述**：配置 DHCP 客户端的默认网关 IP 地址。
- **类型**：字符串（IP 地址）。
- **默认值**：无。
- **说明**：指定网关地址，客户端将使用此地址作为默认网关，通常是路由器的 IP 地址。

## 使用示例

假设你有一个局域网，网段为 `192.168.1.0/24`，你希望为一些设备分配固定的 IP 地址，同时为其它设备动态分配 IP 地址。你可以使用如下配置：

```ini
[plugin.microdhcp]
enable = true
listen_ip = 0.0.0.0
listen_port = 67
static_address_mapping = "00:11:22:33:44:55=192.168.1.100, 00:11:22:33:44:56=192.168.1.101"
dhcp_range_start = 192.168.1.200
dhcp_range_end = 192.168.1.250
dhcp_lease_time = 86400  # 24 小时
dns_servers = "8.8.8.8, 8.8.4.4"
gateway_ip = "192.168.1.1"
```

在此配置中：
- `00:11:22:33:44:55` 和 `00:11:22:33:44:56` 这两台设备将分别获得 IP 地址 `192.168.1.100` 和 `192.168.1.101`。
- 其它设备将从 `192.168.1.200` 到 `192.168.1.250` 范围内动态获取 IP 地址。
- 客户端将使用 `192.168.1.1` 作为网关，DNS 服务器为 `8.8.8.8` 和 `8.8.4.4`。
