# RHILEX 工业边缘网关系统
![image](https://github.com/user-attachments/assets/f02f3900-34a6-4a53-b161-993656e431a1)
<img width="1912" height="920" alt="21" src="https://github.com/user-attachments/assets/f940f195-0695-4915-b987-fc80d05358af" />
<img width="1912" height="920" alt="20" src="https://github.com/user-attachments/assets/c5e92a33-35ed-485e-b28d-60e4400dc108" />
<img width="1912" height="920" alt="19" src="https://github.com/user-attachments/assets/63f39308-726b-4e42-ba79-cf21f0ea27ed" />
<img width="1912" height="920" alt="18" src="https://github.com/user-attachments/assets/1a0ffc7f-9b1e-42c5-b1c9-199298ff95e1" />
<img width="1912" height="920" alt="17" src="https://github.com/user-attachments/assets/a954b156-3803-4541-8474-2a02395876f9" />
<img width="1912" height="920" alt="16" src="https://github.com/user-attachments/assets/2be5d76c-bf9c-40a3-aa74-c4c4ccd3bfaf" />
<img width="1912" height="920" alt="15" src="https://github.com/user-attachments/assets/4ce404b0-b08f-4343-bc56-9b84a59bd354" />
<img width="1912" height="920" alt="14" src="https://github.com/user-attachments/assets/d49c7bc3-1373-48e9-84f6-cc064a6e64c2" />
<img width="1912" height="920" alt="13" src="https://github.com/user-attachments/assets/f21145bd-0b01-4064-91c4-cd51a5819eec" />
<img width="1912" height="920" alt="12" src="https://github.com/user-attachments/assets/b8d748db-bf6d-4384-b763-8daf4072f991" />
<img width="1912" height="920" alt="11" src="https://github.com/user-attachments/assets/3ab5709d-8077-44fd-948d-497e55789997" />
<img width="1912" height="920" alt="10" src="https://github.com/user-attachments/assets/dd0c7800-9f6b-4070-8846-4bc5f1d0b023" />
<img width="1912" height="920" alt="9" src="https://github.com/user-attachments/assets/160d6f94-997c-436e-9e1b-ce560db4a1ab" />
<img width="1912" height="920" alt="8" src="https://github.com/user-attachments/assets/eca2c3dd-1cb5-457d-8b71-8b7785a0a092" />
<img width="1912" height="920" alt="7" src="https://github.com/user-attachments/assets/91f8bc38-112e-4c4d-9579-c63281717f4a" />
<img width="1912" height="920" alt="6" src="https://github.com/user-attachments/assets/aa70a831-5642-478b-a237-1a865ba9fbf1" />
<img width="1912" height="920" alt="5" src="https://github.com/user-attachments/assets/fa97b086-3903-4ebe-9c23-af452498dbfe" />
<img width="1912" height="920" alt="4" src="https://github.com/user-attachments/assets/dd50af71-485f-4aca-9847-ee0a232ffbad" />
<img width="1912" height="920" alt="3" src="https://github.com/user-attachments/assets/966ce3fe-2d07-4faa-9f94-adab6bcf6db2" />
<img width="1912" height="920" alt="2" src="https://github.com/user-attachments/assets/56399e16-ccae-44f7-a0e2-35670a78d548" />
<img width="1912" height="920" alt="1" src="https://github.com/user-attachments/assets/b23b1bbb-b68d-4436-9bcf-ae7052143a0a" />

## 一、项目概述
RHILEX 是一款功能强大的工业边缘网关系统，旨在为工业自动化、物联网等领域提供全面的设备接入、数据处理、协议转换、云边协同等服务。本项目已开源，希望通过社区的力量进一步优化和拓展系统功能，为更多用户提供便捷高效的工业数据处理解决方案。

## 二、功能特点
1. **设备接入与管理**
    - 支持多种工业和物联网协议，包括但不限于：
        - **南向协议**：Modbus（主机/从机模式）、西门子 S7 系列 PLC 采集、SNMP、Bacnet（主/从模式）、HTTP 采集、DLT645 电表协议、CJT188 仪表协议、SZY206 水资源检测协议等，还支持自定义串口协议接入设备。
        - **北向协议**：可以将数据推送到 MQTT Broker、UDP Server、TCP Server、HTTP Server、MongoDB、TdEngine、串口、Semtech UDP Forwarder 及 GreptimeDb 等。
    - 方便的设备管理功能，可对通用串口读写设备、西门子 PLC、各种 Modbus 设备、SNMP 设备、Bacnet 设备、HTTP 数据采集设备、腾讯云物联网平台设备等进行配置、监控和管理。用户可根据需要设置采集频率、寄存器地址、功能码等参数。

2. **数据处理与存储**
    - **数据中心**：存储采集的数据，支持使用 Lua 脚本写入数据，同时提供根据设定条件读取、导出和清空数据的功能。
    - **数据模型**：允许用户根据需求构建数据模型及相应的存储仓库，可定义字段属性，如名称、类型、单位、范围、权限等。数据模型发布后可进行读写操作，类似于数据库的建表过程。
    - **规则引擎**：使用 Lua 脚本编写规则，实现数据的过滤、转换、计算等处理逻辑，并可根据规则触发相应动作，例如根据设定的阈值判断设备状态并执行通知或控制指令。

3. **系统配置与管理**
    - **配置指南**：涵盖了多个方面的配置参数，包括应用程序（日志输出、调试模式、资源限制等）、插件（HTTP API、USB 监控、Modbus 工具等各种插件的启用和参数设置）、传输（串口通信参数）等。
    - **证书管理**：提供证书申请、配置（指定证书路径）、验证（检查证书有效性）的完整流程。
    - **系统设置**：提供查看系统资源、网络状态、设置端口、网卡、路由、WIFI、4G 网络、时间、固件升级、数据备份、用户信息等功能，方便用户对系统进行全面管理。

4. **云边协同**：可与联犀平台协同工作，实现在联犀平台创建产品和设备，并在 RHILEX 中进行接入和配置协同，确保正确的产品-设备-秘钥三元组设置。开启后，支持数据映射与交互，例如将 Modbus 设备数据上传至联犀平台，并能接收云端指令。

5. **二次开发**：基于 RHILEX 框架（采用 AGPL 协议开源）开发，为开发者提供了完善的开发环境搭建指南，包括推荐在 Linux 下开发，详细的启动程序方法，以及开发工具（如 Visual Studio Code 和 Jetbrain Goland）的配置建议。同时，对关键接口（南向、北向、设备、插件接口及其方法定义）进行了详细说明，并提供了设备、北向、插件开发案例及综合案例，方便开发者快速上手。

6. **辅助功能**
    - **轻量应用**：基于 Lua 的扩展脚本系统，提供了如数据转发 MQTT、GPIO 控制等示例，开发人员可根据具体需求灵活开发新应用。
    - **增强插件**：包括 API Server、CRC 计算器、USB 监控器、ICMP 测速、Modbus 扫描、Ngrok 客户端及 Ngrok 内网透传插件等，可有效扩展系统功能。例如，Ngrok 插件可实现将本地端口映射到公网，方便远程访问和调试。
    - **数据遥测**：配备遥测插件，可收集设备运行状态、性能指标等信息，其数据格式公开，用户可通过配置文件开关控制，并遵循数据遥测协议，保障用户权益。
    - **通信模组**：支持多种通信模块（如 Lora、蓝牙、WIFI 等）的接入，提供环境参数及交互流程示例，方便用户集成使用。


## 三、安装与部署
- 前端：https://github.com/hootrhino/rhilex-web 下载前端项目。
- 其他：https://www.hootrhino.com 获取资料。（近期服务器到期，决定不再续费，如果有定制化需求请联系我）


## 四、使用方法

### 设备接入与数据采集
1. 根据设备的协议类型，在设备接入模块中添加新设备并配置相应参数，例如对于串口设备，需配置串口参数、设备地址、功能码等。
2. 导入点位表，您可以手动或批量配置点位，点位表中应包含传感器别名、功能码、采样频率、数据类型等信息。完成配置后，系统将自动开始采集数据，并存储至数据中心。

## 五、贡献指南
我们欢迎社区成员对 RHILEX 项目进行贡献，以下是一些贡献的方式和建议：
1. **代码贡献**：
    - 请先将本项目 fork 到您的 GitHub 账户，在您的分支上进行开发。
    - 确保您的代码遵循项目的编码规范和风格。
    - 提交代码前，运行测试用例，确保代码的正确性和稳定性。
    - 发起 pull request，详细描述您的代码修改内容和目的，等待审核和合并。
2. **文档贡献**：
    - 发现文档中的错误或不足，可以直接修改并提交 pull request。
    - 为项目添加新的使用案例、教程或说明，帮助更多用户更好地使用 RHILEX。
3. **问题反馈**：
    - 若在使用过程中遇到问题，请在 Issues 页面创建新的问题，详细描述问题的症状、操作步骤和系统环境，以便我们快速定位和解决问题。


## 六、许可证
本项目采用 [AGPL2] 许可证。这意味着您可以自由使用、修改和分发本项目，但需要遵守相应的开源协议，请在使用前仔细阅读许可证内容。


## 七、致谢
感谢所有参与 RHILEX 项目开发、测试、使用和推广的人员，是大家的共同努力让这个项目不断发展壮大。


## 八、联系我们
如果您对项目有任何疑问、建议或合作意向，请通过以下方式联系我们：
- ![image](https://github.com/user-attachments/assets/4fb3107e-5307-469e-af9f-0a1a8814eb35)

> RHILEX全部商业仓库代码均公开，未来不做免费答疑。感兴趣请自行研究源码，结构简单注释也有很多。
