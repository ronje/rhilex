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

package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/crontask"
	"github.com/hootrhino/rhilex/typex"
	"github.com/robfig/cron/v3"
)

/*
*
* 初始化路由
*
 */
func InitCronRebootRoute() {
	route := server.RouteGroup(server.ContextUrl("/cronReboot"))
	route.GET("/config", server.AddRoute(GetCronRebootConfig))
	route.POST("/update", server.AddRoute(SetCronRebootConfig))

}

type CronRebootConfigVo struct {
	Enable   *bool  `json:"enable"`
	CronExpr string `json:"cron_expr"`
}

/**
 * 定时重启
 *
 */
func GetCronRebootConfig(c *gin.Context, ruleEngine typex.Rhilex) {
	Config, err := service.GetCronRebootConfig()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(CronRebootConfigVo{
		Enable:   Config.Enable,
		CronExpr: Config.CronExpr,
	}))
}

/**
 * 更新
 *
 */
func SetCronRebootConfig(c *gin.Context, ruleEngine typex.Rhilex) {
	vo := CronRebootConfigVo{}
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, errParse := specParser.Parse(vo.CronExpr); errParse != nil {
		c.JSON(common.HTTP_OK, common.Error400(errParse))
		return
	}
	err1 := service.UpdateMCronRebootConfig(&model.MCronRebootConfig{
		Enable:   vo.Enable,
		CronExpr: vo.CronExpr,
	})
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	if !*vo.Enable {
		errParse := crontask.StopCronRebootCron(vo.CronExpr)
		if errParse != nil {
			c.JSON(common.HTTP_OK, common.Error400(errParse))
			return
		}
	}
	if *vo.Enable {
		errParse := crontask.StartCronRebootCron(vo.CronExpr)
		if errParse != nil {
			c.JSON(common.HTTP_OK, common.Error400(errParse))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
