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

# SMTP 邮件发送器
## 配置

```json
{
  "server": "smtp.example.com",
  "port": 587,
  "user": "your-email@example.com",
  "password": "your-password",
  "subject": "Hello",
  "from": "sender@example.com",
  "to": "recipient@example.com"
}
```

## 示例
```lua
mail:SendSimple('$uuid', 'title', 'Hello world')
```