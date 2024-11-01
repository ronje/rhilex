package main

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

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	archsupport "github.com/hootrhino/rhilex/archsupport"
	"github.com/hootrhino/rhilex/component/activation"
	"github.com/hootrhino/rhilex/engine"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/urfave/cli/v2"
)

func init() {

	go func() {
		for {
			select {
			case <-context.Background().Done():
				return
			default:
				time.Sleep(30 * time.Second)
				runtime.GC()
			}
		}
	}()
	env := os.Getenv("ARCHSUPPORT")
	typex.DefaultVersionInfo.Product = archsupport.CheckVendor(env)
	dist, err := utils.GetOSDistribution()
	if err != nil {
		utils.CLog("Failed to Get OS Distribution:%s", err)
		os.Exit(1)
	}
	typex.DefaultVersionInfo.Dist = dist
	arch := fmt.Sprintf("%s-%s", typex.DefaultVersionInfo.Dist, runtime.GOARCH)
	typex.DefaultVersionInfo.Arch = arch
}

//
//go:generate bash ./gen_info.sh
func main() {
	app := &cli.App{
		Name:  "RHILEX STREAM SYSTEM",
		Usage: "For more, please refer to: https://www.hootrhino.com",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start rhilex with config: -config=/path/rhilex.ini",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db",
						Usage: "rhilex database",
						Value: "rhilex.db",
					},
					&cli.StringFlag{
						Name:  "config",
						Usage: "rhilex config",
						Value: "rhilex.ini",
					},
				},
				Action: func(c *cli.Context) error {
					utils.CLog(typex.Banner)
					utils.ShowGGpuAndCpuInfo()
					utils.ShowIpAddress()
					pid := os.Getpid()
					err := os.WriteFile(ossupport.MainExePidPath, []byte(fmt.Sprintf("%d", pid)), 0755)
					if err != nil {
						utils.CLog("[RHILEX RUN] Write Pid File Failed:%s", err)
						return nil
					}
					engine.RunRhilex(c.String("config"))
					if utils.PathExists(ossupport.MainExePidPath) {
						os.Remove(ossupport.MainExePidPath)
					}
					utils.CLog("[RHILEX RUN] Stop rhilex successfully.")
					return nil
				},
			},
			{
				Name:   "upgrade",
				Hidden: true,
				Usage:  "! JUST FOR Upgrade FirmWare",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "upgrade",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "inipath",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "licpath",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "keypath",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "rundbpath",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: "",
					},
				},
				Action: func(c *cli.Context) error {
					flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
					file, err := os.OpenFile(ossupport.UpgradeLogPath, flag, 0755)
					if err != nil {
						utils.CLog(err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						utils.CLog("[RHILEX UPGRADE] Write Upgrade Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							utils.CLog("[RHILEX UPGRADE] Remove Upgrade Lock File error:%s", err.Error())
							return
						}
						utils.CLog("[RHILEX UPGRADE] Remove Upgrade Lock File Finished")
					}()
					if !c.Bool("upgrade") {
						utils.CLog("[RHILEX UPGRADE] Nothing todo")
						return nil
					}
					utils.CLog("[RHILEX BACKUP] Start backup ")
					var errBob error
					errBob = ossupport.BackupOldVersion(ossupport.MainWorkDir+ossupport.GetExePath(),
						ossupport.OldBackupDir+ossupport.GetExePath())
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup rhilex Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("inipath"), ossupport.OldBackupDir+"rhilex.ini")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup rhilex.ini Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("rundbpath"), ossupport.OldBackupDir+"rhilex.db")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup rhilex.db Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("keypath"), ossupport.OldBackupDir+"license.key")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup License.Key Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("licpath"), ossupport.OldBackupDir+"license.lic")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup License.Lic Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.DataCenterPath, ossupport.OldBackupDir+"rhilex_datacenter.db")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup DataCenter Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.LostCacheDataPath, ossupport.OldBackupDir+"rhilex_lostcache.db")
					if errBob != nil {
						utils.CLog("[RHILEX BACKUP] Backup LostCacheData Failed: %s", errBob)
						return nil
					}
					utils.CLog("[RHILEX BACKUP] Backup finished")
					utils.CLog("[RHILEX UPGRADE] Unzip Firmware")
					cwd, err := os.Getwd()
					if err != nil {
						utils.CLog("[RHILEX UPGRADE] Getwd error: %v", err)
						return nil
					}
					if err := ossupport.UnzipFirmware(ossupport.FirmwarePath, cwd); err != nil {
						utils.CLog("[RHILEX UPGRADE] Unzip Firmware error:%s", err.Error())
						return nil
					}
					utils.CLog("[RHILEX UPGRADE] Unzip Firmware finished")
					utils.CLog("[RHILEX UPGRADE] Upgrade finished")
					return nil
				},
			},
			// 回滚 TODO
			{
				Name:   "rollback",
				Usage:  "! JUST FOR Rollback ",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "rollback",
						Usage: "! THIS PARAMENT IS JUST FOR Rollback ",
						Value: false,
					},
				},
				Action: func(c *cli.Context) error {
					utils.CLog("[RHILEX ROLLBACK] Rollback Process Started")
					var errBob error
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+ossupport.GetExePath(), ossupport.GetExePath())
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex.ini", ossupport.RunIniPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex.db", ossupport.RunDbPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"license.key", ossupport.LicenseKeyPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"license.lic", ossupport.LicenseLicPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex_datacenter.db", ossupport.DataCenterPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex_lostcache.db", ossupport.LostCacheDataPath)
					if errBob != nil {
						utils.CLog("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					utils.CLog("[RHILEX ROLLBACK] Rollback Process Exited")
					return nil
				},
			},
			// 数据恢复
			{
				Name:   "recover",
				Usage:  "! JUST FOR Recover Data",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "recover",
						Usage: "! THIS PARAMENT IS JUST FOR Recover Data",
						Value: false,
					},
				},
				Action: func(c *cli.Context) error {
					file, err := os.Create(ossupport.RecoverLogPath)
					if err != nil {
						utils.CLog(err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						utils.CLog("[DATA RECOVER] Write Recover Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							utils.CLog("[DATA RECOVER] Remove Recover Lock File error:%s", err.Error())
							return
						}
						utils.CLog("[DATA RECOVER] Remove Recover Lock File Finished")
					}()
					if runtime.GOOS != "linux" {
						utils.CLog("[DATA RECOVER] Only Support Linux")
						return nil
					}

					if !c.Bool("recover") {
						utils.CLog("[DATA RECOVER] Nothing todo")
						return nil
					}
					utils.CLog("[DATA RECOVER] Remove Old Db File")
					if err := os.Remove(ossupport.RunDbPath); err != nil {
						utils.CLog("[DATA RECOVER] Remove Main COnfig Db error:%s", err.Error())
						return nil
					}
					if err := os.Remove(ossupport.DataCenterPath); err != nil {
						utils.CLog("[DATA RECOVER] Remove Data Center Db error:%s", err.Error())
						return nil
					}
					if err := os.Remove(ossupport.LostCacheDataPath); err != nil {
						utils.CLog("[DATA RECOVER] Remove Lost Cache Db error:%s", err.Error())
						return nil
					}
					utils.CLog("[DATA RECOVER] Remove Old Db File Finished")
					utils.CLog("[DATA RECOVER] Move New Db File")
					if err := ossupport.MoveFile(ossupport.RecoveryDbPath,
						ossupport.RunDbPath); err != nil {
						utils.CLog("[DATA RECOVER] Move New Db File error:%s", err.Error())
						return nil
					}
					if err := ossupport.MoveFile(ossupport.RecoveryDataCenterPath,
						ossupport.DataCenterPath); err != nil {
						utils.CLog("[DATA RECOVER] Move DataCenter File error:%s", err.Error())
						return nil
					}
					utils.CLog("[DATA RECOVER] Move New Db File Finished")
					utils.CLog("[DATA RECOVER] Try to Restart rhilex")
					if err := ossupport.RestartRhilex(); err != nil {
						utils.CLog("[DATA RECOVER] Restart rhilex error:%s", err.Error())
					} else {
						utils.CLog("[DATA RECOVER] Restart rhilex success, Recover Process Exited")
					}
					os.Exit(0)
					return nil
				},
			},
			{
				Name:   "active",
				Usage:  "active -H ${HOST} --SN ${SN} -U ${Username} -P ${Password} -IF ${Iface name} -MAC ${MAC Address}",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "H",
						Usage: "server ip",
					},
					&cli.StringFlag{
						Name:  "SN",
						Usage: "serial number",
					},
					&cli.StringFlag{
						Name:  "U",
						Usage: "admin username",
					},
					&cli.StringFlag{
						Name:  "P",
						Usage: "admin password",
					},
					&cli.StringFlag{
						Name:  "IF",
						Usage: "interface name",
					},
					&cli.StringFlag{
						Name:  "MAC",
						Usage: "mac address",
					},
				},

				Action: func(c *cli.Context) error {
					host := c.String("H")
					if host == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing 'host' parameter")
					}
					sn := c.String("SN")
					if sn == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing 'SN' parameter")
					}
					iface := c.String("IF")
					if iface == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin 'IF' parameter")
					}
					mac := c.String("MAC")
					if iface == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin 'MAC' parameter")
					}
					U := c.String("U")
					if U == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin 'U' parameter")
					}
					P := c.String("P")
					if P == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin 'P' parameter")
					}

					key, lic, errGetLicense := activation.GetLicense(host, sn, iface, mac, U, P)
					if errGetLicense != nil {
						utils.CLog("[LICENCE ACTIVE]: Get License failed:%s", errGetLicense)
						return nil
					}
					licFile, err1 := os.Create("./license.lic")
					if err1 != nil {
						utils.CLog("[LICENCE ACTIVE]: Create License failed:%s", err1)
						return nil
					}
					licFile.Write([]byte(lic))
					keyFile, err2 := os.Create("./license.key")
					if err2 != nil {
						utils.CLog("[LICENCE ACTIVE]: Create Key failed:%s", err2)
						return nil
					}
					keyFile.Write([]byte(key))
					utils.CLog("[LICENCE ACTIVE]: Get License success. save to ./license.key, ./license.lic")
					defer licFile.Close()
					defer keyFile.Close()
					return nil
				},
			},
			{
				Name:   "validate",
				Usage:  "validate rhilex license",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key",
						Usage: "license key path",
					},
					&cli.StringFlag{
						Name:  "lic",
						Usage: "license path",
					},
				},
				Action: func(c *cli.Context) error {
					keyPath := c.String("key")
					if keyPath == "" {
						utils.CLog("[LICENCE ACTIVE]: missing admin 'key' parameter")
						return nil
					}
					licPath := c.String("lic")
					if licPath == "" {
						utils.CLog("[LICENCE ACTIVE]: missing admin 'lic' parameter")
						return nil
					}
					// rhilex validate -lic ./license.lic -key ./license.key
					LocalLicense, err := utils.ValidateLicense(keyPath, licPath)
					if err != nil {
						utils.CLog("[LICENCE ACTIVE]: Validate License Failed: %s", err.Error())
						return nil
					}
					fmt.Println(LocalLicense.ToString())
					return nil
				},
			},
			// version
			{
				Name:  "version",
				Usage: "Show rhilex Current Version",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "version",
						Usage: "rhilex version",
					},
				},
				Action: func(*cli.Context) error {
					version := fmt.Sprintf("[%v-%v-%v]",
						runtime.GOOS, runtime.GOARCH, typex.MainVersion)
					utils.CLog("[*] Version: " + version)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
