package apis

import (
	"regexp"
	"runtime"
	"sync"

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
  - 设置时间、时区
  - sudo date -s "2023-08-07 15:30:00"
    获取时间: date "+%Y-%m-%d %H:%M:%S" -> 2023-08-07 15:30:00
*/
func SetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	DtoCfg := TimeVo{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if validTimeZone(DtoCfg.SysTimeZone) {
		c.JSON(common.HTTP_OK, common.Error("Invalid TimeZone:"+DtoCfg.SysTimeZone))
		return
	}

	err1 := ossupport.SetSystemTime(DtoCfg.SysTime)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	err2 := ossupport.SetTimeZone(DtoCfg.SysTimeZone)
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 获取系统时间
*
 */
func GetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
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

var (
	once     sync.Once
	timeZone *regexp.Regexp
)

func validTimeZone(timezone string) bool {
	once.Do(func() {
		regexPattern := `^[A-Za-z]+/[A-Za-z_]+$`
		timeZone = regexp.MustCompile(regexPattern)
	})

	return timeZone.MatchString(timezone)
}

/*
*
* 更新时间
*
 */
func UpdateTimeByNtp(c *gin.Context, ruleEngine typex.Rhilex) {
	if err := ossupport.UpdateTimeByNtp(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}
