// Copyright (C) 2024 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package crontask

import (
	"context"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/component/shellengine"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/robfig/cron/v3"
)

var __DefaultCronRebootExecutor *CronRebootExecutor

type CronRebootExecutor struct {
	LinuxBashShell *shellengine.LinuxBashShell
	Cron           *cron.Cron
	CronEntryID    cron.EntryID
}
type MCronRebootConfig struct {
	ID        uint
	CreatedAt time.Time
	Enable    bool
	CronExpr  string
}

func InitCronRebootExecutor(rhilex typex.Rhilex) {
	__DefaultCronRebootExecutor = &CronRebootExecutor{
		LinuxBashShell: shellengine.InitLinuxBashShell(rhilex),
		Cron:           cron.New(cron.WithChain()),
		CronEntryID:    -100,
	}
	m := new(MCronRebootConfig)
	m.ID = 1
	m.CreatedAt = time.Now()
	m.Enable = false
	m.CronExpr = "0 0 0 0 0"
	err := interdb.DB().Model(m).FirstOrCreate(m).Error
	if err != nil {
		glogger.GLogger.Error(err)
		return
	}
	if m.Enable {
		errParse := StartCronRebootCron(m.CronExpr)
		if errParse != nil {
			glogger.GLogger.Error(errParse)
			return
		}
	}
}

/**
 * StopCronRebootCron
 *
 */
func StopCronRebootCron(expr string) error {
	if __DefaultCronRebootExecutor.CronEntryID > 0 {
		glogger.GLogger.Info("Remove Cron Reboot Cron:", expr)
		__DefaultCronRebootExecutor.Cron.Remove(__DefaultCronRebootExecutor.CronEntryID)
	}
	__DefaultCronRebootExecutor.Cron.Stop()
	return nil
}

/**
 * StartCronRebootCron
 *
 */
func StartCronRebootCron(expr string) error {
	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, errParse := specParser.Parse(expr); errParse != nil {
		return errParse
	}
	__DefaultCronRebootExecutor.Cron.Remove(__DefaultCronRebootExecutor.CronEntryID)
	__DefaultCronRebootExecutor.Cron.Stop()
	var err error
	__DefaultCronRebootExecutor.CronEntryID, err = __DefaultCronRebootExecutor.Cron.AddFunc(expr, func() {
		if core.GlobalConfig.AppDebugMode {
			glogger.GLogger.Debug("Start Cron Reboot Cron:", expr)
		}
		if runtime.GOOS == "linux" {
			__DefaultCronRebootExecutor.LinuxBashShell.JustRun(context.Background(), "reboot")
		}
		if runtime.GOOS == "windows" {
			glogger.GLogger.Error("Windows Not Support Reboot, Please Set On Windows Control Panel:", expr)
		}
	})
	if err != nil {
		return err
	}
	__DefaultCronRebootExecutor.Cron.Start()
	return nil
}

/**
 * Stop
 *
 */
func StopCronRebootExecutor() {
	__DefaultCronRebootExecutor.Cron.Remove(__DefaultCronRebootExecutor.CronEntryID)
	__DefaultCronRebootExecutor.Cron.Stop()
}
