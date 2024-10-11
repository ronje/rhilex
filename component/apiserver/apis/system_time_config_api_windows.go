package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 时区设置
*
 */
type TimeVo struct {
	SysTime     string `json:"sysTime"`
	SysTimeZone string `json:"sysTimeZone"`
	EnableNtp   bool   `json:"enableNtp"`
}

/*
*
* 获取系统时间
*
 */
func GetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	SysTime, err := ossupport.GetSystemTime()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	SysTimeZone, err := ossupport.GetTimeZone()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.OkWithData(TimeVo{
		EnableNtp:   true,
		SysTime:     SysTime,
		SysTimeZone: SysTimeZone.CurrentTimezone,
	}))
}

/*
*
* 更新时间
*
 */
func UpdateTimeByNtp(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/**
 * 配置时间
 *
 */
func SetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.Ok())
}
