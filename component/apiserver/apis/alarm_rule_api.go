package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

func LoadAlarmRuleRoute() {
	api := server.RouteGroup(server.ContextUrl("/alarm_rule"))
	api.GET(("/list"), server.AddRoute(AlarmRuleList))
	api.POST(("/create"), server.AddRoute(CreateAlarmRule))
	api.PUT(("/update"), server.AddRoute(UpdateAlarmRule))
	api.DELETE(("/del"), server.AddRoute(DeleteAlarmRule))
	api.GET(("/detail"), server.AddRoute(AlarmRuleDetail))
}

// 告警规则
type AlarmRuleVo struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Expr        string `json:"expr"`
	Interval    uint64 `json:"interval"`
	Description string `json:"description"`
}

/*
*
* AlarmRule
*
 */
func AlarmRuleList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count, AlarmRules := service.PageAlarmRule(pager.Current, pager.Size)
	AlarmRuleVos := []AlarmRuleVo{}
	for _, AlarmRule := range AlarmRules {
		AlarmRuleVos = append(AlarmRuleVos, AlarmRuleVo{
			UUID:        AlarmRule.UUID,
			Name:        AlarmRule.Name,
			Expr:        AlarmRule.Expr,
			Interval:    AlarmRule.Interval,
			Description: AlarmRule.Description,
		})
	}
	Result := service.WrapPageResult(*pager, AlarmRuleVos, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 详情
func AlarmRuleDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	AlarmRule, err1 := service.GetMAlarmRuleWithUUID(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err1))
		return
	}
	web_data := AlarmRuleVo{
		UUID:        AlarmRule.UUID,
		Name:        AlarmRule.Name,
		Expr:        AlarmRule.Expr,
		Interval:    AlarmRule.Interval,
		Description: AlarmRule.Description,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(web_data))
}

/*
*
* 删除
*
 */
func DeleteAlarmRule(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	err := service.DeleteAlarmRule(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* Create
*
 */

func CreateAlarmRule(c *gin.Context, ruleEngine typex.Rhilex) {
	AlarmRule := AlarmRuleVo{}
	if err := c.ShouldBindJSON(&AlarmRule); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Model := model.MAlarmRule{
		UUID:        utils.AlarmRuleUuid(),
		Name:        AlarmRule.Name,
		Expr:        AlarmRule.Expr,
		Interval:    AlarmRule.Interval,
		Description: AlarmRule.Description,
	}
	if err := service.InsertAlarmRule(&Model); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新
*
 */
func UpdateAlarmRule(c *gin.Context, ruleEngine typex.Rhilex) {
	AlarmRule := AlarmRuleVo{}
	if err := c.ShouldBindJSON(&AlarmRule); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Model := model.MAlarmRule{
		UUID:        AlarmRule.UUID,
		Name:        AlarmRule.Name,
		Expr:        AlarmRule.Expr,
		Interval:    AlarmRule.Interval,
		Description: AlarmRule.Description,
	}
	if err := service.UpdateAlarmRule(&Model); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
