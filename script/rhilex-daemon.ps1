# 安装路径
$installPath = "C:\Program Files\Rhilex"
# 服务名称
$serviceName = "RhilexService"

# 安装函数
function Install-Rhilex {
    try {
        Write-Host "Installing Rhilex..."

        # 创建安装目录
        if (-not (Test-Path -Path $installPath)) {
            New-Item -ItemType Directory -Path $installPath -Force | Out-Null
        }

        # 复制文件
        Copy-Item -Path ".\rhilex.exe" -Destination $installPath -Force
        Copy-Item -Path ".\license.*" -Destination $installPath -Force
        Copy-Item -Path ".\*.ini" -Destination $installPath -Force

        # 生成系统服务脚本
        $serviceScriptPath = Join-Path -Path $installPath -ChildPath "RhilexService.ps1"
        @"
[CmdletBinding()]
Param()

`$exePath = Join-Path -Path 'C:\Program Files\Rhilex' -ChildPath 'rhilex.exe'
Start-Process -FilePath `$exePath -NoNewWindow
"@ | Out-File -FilePath $serviceScriptPath -Encoding UTF8

        # 创建服务
        $serviceBinaryPath = "powershell.exe -WindowStyle Hidden -File $serviceScriptPath"
        New-Service -Name $serviceName -BinaryPathName $serviceBinaryPath -DisplayName "Rhilex Service" -StartupType Automatic

        # 启动服务
        Start-Service -Name $serviceName

        Write-Host "Rhilex installed successfully."
    }
    catch {
        Write-Error "An error occurred during installation: $_"
    }
}

# 卸载函数
function Uninstall-Rhilex {
    try {
        Write-Host "Uninstalling Rhilex..."

        # 停止并删除服务
        if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
            Stop-Service -Name $serviceName -Force
            Remove-Service -Name $serviceName -Force
        }

        # 删除安装目录及文件
        if (Test-Path -Path $installPath) {
            Remove-Item -Path $installPath -Recurse -Force
        }

        Write-Host "Rhilex uninstalled successfully."
    }
    catch {
        Write-Error "An error occurred during uninstallation: $_"
    }
}

# 帮助函数
function Show-Help {
    Write-Host "Usage: powershell -File $($MyInvocation.MyCommand.Name) [action]"
    Write-Host ""
    Write-Host "Actions:"
    Write-Host "  install    Install Rhilex. This will copy rhilex.exe, license.*, and *.ini files to the installation directory,"
    Write-Host "              create a system service script, and start the service."
    Write-Host "  uninstall  Uninstall Rhilex. This will stop and remove the system service, and delete all related files."
    Write-Host "  help       Display this help message."
}

# 主脚本逻辑
if ($args[0] -eq "install") {
    Install-Rhilex
} elseif ($args[0] -eq "uninstall") {
    Uninstall-Rhilex
} elseif ($args[0] -eq "help") {
    Show-Help
} else {
    Write-Host "Invalid action. Please specify 'install', 'uninstall', or 'help'."
}