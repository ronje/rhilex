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


# Rhilex服务操作指南
本指南将帮助您在Windows操作系统中安装、卸载和查看Rhilex服务的状态。请按照以下步骤操作：
## 前提条件
- 确保您已以管理员权限登录到Windows。
- 确认`rhilex.exe`和`rhilex.ini`文件位于同一个目录中。
- 下载并保存提供的PowerShell脚本`Install-RhilexService.ps1`到`rhilex.exe`所在的目录。
## 安装Rhilex服务
1. **打开命令提示符（管理员）**：
   - 按下`Win + R`键，输入`cmd`，然后按`Enter`。
   - 右键点击命令提示符窗口，选择“以管理员身份运行”。
2. **导航到脚本所在目录**：
   - 在命令提示符中，输入以下命令并按`Enter`：
     ```
     cd 路径\到\rhilex.exe\的目录
     ```
     请将`路径\到\rhilex.exe\的目录`替换为实际的路径。
3. **运行安装脚本**：
   - 输入以下命令并按`Enter`：
     ```
     .\Install-RhilexService.ps1 -action install
     ```
   - 如果提示权限问题，请确认是否以管理员身份运行命令提示符。
4. **确认服务安装成功**：
   - 安装完成后，您可以通过服务管理器检查Rhilex服务是否已成功安装并设置为自动启动。
## 卸载Rhilex服务
1. **打开命令提示符（管理员）**：
   - 按照上述步骤打开命令提示符。
2. **导航到脚本所在目录**：
   - 使用`cd`命令导航到包含`rhilex.exe`的目录。
3. **运行卸载脚本**：
   - 输入以下命令并按`Enter`：
     ```
     .\Install-RhilexService.ps1 -action uninstall
     ```
4. **确认服务卸载成功**：
   - 卸载完成后，可以通过服务管理器检查Rhilex服务是否已被移除。
## 查看Rhilex服务状态
1. **打开命令提示符（管理员）**：
   - 按照上述步骤打开命令提示符。
2. **导航到脚本所在目录**：
   - 使用`cd`命令导航到包含`rhilex.exe`的目录。
3. **运行查看状态脚本**：
   - 输入以下命令并按`Enter`：
     ```
     .\Install-RhilexService.ps1 -action status
     ```
   - 脚本将显示Rhilex服务的当前状态。
## 注意事项
- 在执行任何操作之前，请确保已保存所有重要数据。
- 如果服务在卸载后仍然存在，您可能需要手动通过服务管理器进行移除。
- 如果遇到任何问题，请检查脚本是否有正确的执行权限，或者联系技术支持。
