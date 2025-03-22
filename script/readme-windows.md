# Rhilex 安装卸载脚本说明文档

## 1. 脚本概述
该 PowerShell 脚本用于实现 Rhilex 应用程序的安装、卸载功能，并提供帮助信息。脚本会将 `rhilex.exe`、`license.*`、`*.ini` 这些文件复制到指定的安装目录，同时生成系统服务脚本，在卸载时会删除所有相关文件。

## 2. 脚本使用前提
- 运行此脚本需要具备管理员权限，因为脚本会涉及到创建和删除 Windows 服务，以及操作 `C:\Program Files` 目录下的文件。
- 脚本运行前，请确保 `rhilex.exe`、`license.*`、`*.ini` 这些文件存在于脚本所在的同一目录下。

## 3. 脚本参数及功能说明

### 3.1 安装（`install`）
- **参数**：`install`
- **功能**：
  - **创建安装目录**：如果 `C:\Program Files\Rhilex` 目录不存在，脚本会自动创建该目录。
  - **复制文件**：将脚本所在目录下的 `rhilex.exe`、`license.*` 和 `*.ini` 文件复制到 `C:\Program Files\Rhilex` 目录。
  - **生成系统服务脚本**：在安装目录下生成一个名为 `RhilexService.ps1` 的 PowerShell 脚本，该脚本的作用是启动 `rhilex.exe`。
  - **创建并启动服务**：使用 `New-Service` 命令创建一个名为 `RhilexService` 的 Windows 服务，将其设置为自动启动，并立即启动该服务。

### 3.2 卸载（`uninstall`）
- **参数**：`uninstall`
- **功能**：
  - **停止并删除服务**：如果 `RhilexService` 服务存在，脚本会先停止该服务，然后将其从系统服务列表中删除。
  - **删除安装目录及文件**：删除 `C:\Program Files\Rhilex` 目录及其下的所有文件和子目录。

### 3.3 帮助（`help`）
- **参数**：`help`
- **功能**：显示脚本的使用说明，包括脚本的基本用法和各个操作的功能描述。

## 4. 脚本使用示例

### 4.1 安装 Rhilex
```powershell
powershell -File InstallUninstallScript.ps1 install
```

### 4.2 卸载 Rhilex
```powershell
powershell -File InstallUninstallScript.ps1 uninstall
```

### 4.3 查看帮助信息
```powershell
powershell -File InstallUninstallScript.ps1 help
```

## 5. 注意事项
- 脚本运行过程中如果出现错误，会输出相应的错误信息，方便用户进行排查。
- 请确保在运行脚本前关闭所有与 Rhilex 相关的程序和服务，以免影响安装或卸载过程。
- 若要再次安装 Rhilex，建议先卸载已有的版本，以避免文件冲突和服务冲突。