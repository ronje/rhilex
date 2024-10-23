// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ossupport

/*
*
* Linux系统下的一些和应用交互的系统级路径
*
 */
const (
	// rhilex 工作目录
	MainWorkDir = "/usr/local/rhilex/"
	// RHILEX Main
	MainExePath = MainWorkDir + "rhilex"
	// RHILEX Config
	RunConfigPath = MainWorkDir + "rhilex.ini"
	// RHILEX Database
	RunDbPath = MainWorkDir + "rhilex.db"
	// 证书公钥位置
	LicenseKeyPath = MainWorkDir + "license.key"
	// 证书位置
	LicenseLicPath = MainWorkDir + "license.lic"
	// RHILEX 备份回滚目录
	OldBackupDir = "/usr/local/rhilex/old/"
	// 数据中心
	DataCenterPath = MainWorkDir + "rhilex_datacenter.db"
	// 离线缓存的数据
	LostCacheDataPath = MainWorkDir + "rhilex_lostcache.db"
	// 固件保存路径
	FirmwarePath = MainWorkDir + "upload/Firmware/Firmware.zip"
	// 升级日志
	UpgradeLogPath = MainWorkDir + "rhilex-upgrade-log.txt"
	// 运行时日志
	RunningLogPath = MainWorkDir + "rhilex-running-log.txt"
	// 数据恢复日志
	RecoverLogPath = MainWorkDir + "rhilex-recover-log.txt"
	// 备份锁
	BackupLockPath = "/var/run/rhilex-upgrade.lock"
	// 升级锁
	UpgradeLockPath = BackupLockPath
	// 数据备份
	RecoverBackupPath = MainWorkDir + "upload/Backup/"
	// 备份数据库
	RecoveryDbPath = RecoverBackupPath + "rhilex.db"
	// 数据中心库
	RecoveryDataCenterPath = RecoverBackupPath + "rhilex_datacenter.db"
)
