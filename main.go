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

	"github.com/hootrhino/rhilex/engine"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/periphery"
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
	typex.DefaultVersionInfo.Product = periphery.CheckVendor(env)
	dist, err := utils.GetOSDistribution()
	if err != nil {
		glogger.DefaultOutput("Failed to Get OS Distribution:%s", err)
		os.Exit(1)
	}
	typex.DefaultVersionInfo.Dist = dist
	arch := fmt.Sprintf("%s-%s", typex.DefaultVersionInfo.Dist, runtime.GOARCH)
	typex.DefaultVersionInfo.Arch = arch
}

//
//go:generate bash ./gen_info.sh
func main() {
	defer utils.WritePanicStack()
	app := &cli.App{
		Name:        "rhilex",
		Usage:       "rhilex",
		Description: "RHILEX is a system dedicated to edge gateways.\nMore information visit https://www.hootrhino.com.",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start RHILEX",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db",
						Usage: "specific rhilex database",
						Value: "rhilex.db",
					},
					&cli.StringFlag{
						Name:  "config",
						Usage: "specific rhilex config database",
						Value: "rhilex.ini",
					},
				},
				Action: func(c *cli.Context) error {
					glogger.DefaultOutput("%s", typex.Banner)
					utils.ShowIpAddress()
					pid := os.Getpid()
					err := os.WriteFile(ossupport.MainExePidPath, []byte(fmt.Sprintf("%d", pid)), 0755)
					if err != nil {
						glogger.DefaultOutput("[RHILEX RUN] Write Pid File Failed:%s", err)
						return nil
					}
					engine.RunRhilex(c.String("config"))
					if utils.PathExists(ossupport.MainExePidPath) {
						os.Remove(ossupport.MainExePidPath)
					}
					glogger.DefaultOutput("[RHILEX RUN] Stop rhilex successfully.")
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
						glogger.DefaultOutput("%s", err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						glogger.DefaultOutput("[RHILEX UPGRADE] Write Upgrade Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							glogger.DefaultOutput("[RHILEX UPGRADE] Remove Upgrade Lock File error:%s", err.Error())
							return
						}
						glogger.DefaultOutput("[RHILEX UPGRADE] Remove Upgrade Lock File Finished")
					}()
					if !c.Bool("upgrade") {
						glogger.DefaultOutput("[RHILEX UPGRADE] Nothing todo")
						return nil
					}
					glogger.DefaultOutput("[RHILEX BACKUP] Start backup ")
					var errBob error
					errBob = ossupport.BackupOldVersion(ossupport.MainWorkDir+ossupport.GetExePath(),
						ossupport.OldBackupDir+ossupport.GetExePath())
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup rhilex Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("inipath"), ossupport.OldBackupDir+"rhilex.ini")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup rhilex.ini Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("rundbpath"), ossupport.OldBackupDir+"rhilex.db")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup rhilex.db Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("keypath"), ossupport.OldBackupDir+"license.key")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup License.Key Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(c.String("licpath"), ossupport.OldBackupDir+"license.lic")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup License.Lic Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.DataCenterPath, ossupport.OldBackupDir+"rhilex_datacenter.db")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup DataCenter Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.LostCacheDataPath, ossupport.OldBackupDir+"rhilex_lostcache.db")
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX BACKUP] Backup LostCacheData Failed: %s", errBob)
						return nil
					}
					glogger.DefaultOutput("[RHILEX BACKUP] Backup finished")
					glogger.DefaultOutput("[RHILEX UPGRADE] Unzip Firmware")
					cwd, err := os.Getwd()
					if err != nil {
						glogger.DefaultOutput("[RHILEX UPGRADE] Getwd error: %v", err)
						return nil
					}
					if err := ossupport.UnzipFirmware(ossupport.FirmwarePath, cwd); err != nil {
						glogger.DefaultOutput("[RHILEX UPGRADE] Unzip Firmware error:%s", err.Error())
						return nil
					}
					glogger.DefaultOutput("[RHILEX UPGRADE] Unzip Firmware finished")
					glogger.DefaultOutput("[RHILEX UPGRADE] Upgrade finished")
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
					glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Process Started")
					var errBob error
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+ossupport.GetExePath(), ossupport.GetExePath())
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex.ini", ossupport.RunIniPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex.db", ossupport.RunDbPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"license.key", ossupport.LicenseKeyPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"license.lic", ossupport.LicenseLicPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex_datacenter.db", ossupport.DataCenterPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					errBob = ossupport.BackupOldVersion(ossupport.OldBackupDir+"rhilex_lostcache.db", ossupport.LostCacheDataPath)
					if errBob != nil {
						glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Failed: %s", errBob)
						return nil
					}
					glogger.DefaultOutput("[RHILEX ROLLBACK] Rollback Process Exited")
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
						glogger.DefaultOutput("%s", err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Write Recover Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							glogger.DefaultOutput("[DATA RECOVER] Remove Recover Lock File error:%s", err.Error())
							return
						}
						glogger.DefaultOutput("[DATA RECOVER] Remove Recover Lock File Finished")
					}()
					if runtime.GOOS != "linux" {
						glogger.DefaultOutput("[DATA RECOVER] Only Support Linux")
						return nil
					}

					if !c.Bool("recover") {
						glogger.DefaultOutput("[DATA RECOVER] Nothing todo")
						return nil
					}
					glogger.DefaultOutput("[DATA RECOVER] Remove Old Db File")
					if err := os.Remove(ossupport.RunDbPath); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Remove Main COnfig Db error:%s", err.Error())
						return nil
					}
					if err := os.Remove(ossupport.DataCenterPath); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Remove Data Center Db error:%s", err.Error())
						return nil
					}
					if err := os.Remove(ossupport.LostCacheDataPath); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Remove Lost Cache Db error:%s", err.Error())
						return nil
					}
					glogger.DefaultOutput("[DATA RECOVER] Remove Old Db File Finished")
					glogger.DefaultOutput("[DATA RECOVER] Move New Db File")
					if err := ossupport.MoveFile(ossupport.RecoveryDbPath,
						ossupport.RunDbPath); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Move New Db File error:%s", err.Error())
						return nil
					}
					if err := ossupport.MoveFile(ossupport.RecoveryDataCenterPath,
						ossupport.DataCenterPath); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Move DataCenter File error:%s", err.Error())
						return nil
					}
					glogger.DefaultOutput("[DATA RECOVER] Move New Db File Finished")
					glogger.DefaultOutput("[DATA RECOVER] Try to Restart rhilex")
					if err := ossupport.RestartRhilex(); err != nil {
						glogger.DefaultOutput("[DATA RECOVER] Restart rhilex error:%s", err.Error())
					} else {
						glogger.DefaultOutput("[DATA RECOVER] Restart rhilex success, Recover Process Exited")
					}
					os.Exit(0)
					return nil
				},
			},
			// version
			{
				Name:        "version",
				Usage:       "Show rhilex Current Version",
				Description: "Show rhilex Current Version",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "version",
					},
				},
				Action: func(*cli.Context) error {
					version := fmt.Sprintf("[%v-%v-%v]",
						runtime.GOOS, runtime.GOARCH, typex.MainVersion)
					glogger.DefaultOutput("[*] Version: %s", version)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
