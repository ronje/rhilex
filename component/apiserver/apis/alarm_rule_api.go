package apis

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/component/alarmcenter"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
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
	api.POST(("/testRule"), server.AddRoute(AlarmRuleTest))
}

// 输出规则
type ExprDefineVo struct {
	Expr      string `json:"expr"`
	EventType string `json:"eventType"`
}

// 告警规则
type AlarmRuleVo struct {
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`
	ExprDefine  []ExprDefineVo `json:"exprDefine"`
	Interval    uint64         `json:"interval"`
	Threshold   uint64         `json:"threshold"`
	HandleId    string         `json:"handleId"` // 事件处理器，目前是北向ID
	Description string         `json:"description"`
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
		MOutput := AlarmRule.GetExprDefine()
		exprDef := []ExprDefineVo{}
		for _, output := range MOutput {
			exprDef = append(exprDef, ExprDefineVo{
				Expr:      output.Expr,
				EventType: output.EventType,
			})
		}
		AlarmRuleVos = append(AlarmRuleVos, AlarmRuleVo{
			UUID:        AlarmRule.UUID,
			Name:        AlarmRule.Name,
			ExprDefine:  exprDef,
			Interval:    AlarmRule.Interval,
			Threshold:   AlarmRule.Threshold,
			HandleId:    AlarmRule.HandleId,
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
	MOutput := AlarmRule.GetExprDefine()
	exprDef := []ExprDefineVo{}
	for _, output := range MOutput {
		exprDef = append(exprDef, ExprDefineVo{
			Expr:      output.Expr,
			EventType: output.EventType,
		})
	}
	web_data := AlarmRuleVo{
		UUID:        AlarmRule.UUID,
		Name:        AlarmRule.Name,
		ExprDefine:  exprDef,
		Interval:    AlarmRule.Interval,
		Threshold:   AlarmRule.Threshold,
		HandleId:    AlarmRule.HandleId,
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
	alarmcenter.RemoveExpr(uuid)
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
	for _, Output := range AlarmRule.ExprDefine {
		ok, errVerifyExpr := alarmcenter.VerifyExpr(Output.Expr)
		if errVerifyExpr != nil {
			c.JSON(common.HTTP_OK, common.Error400(errVerifyExpr))
			return
		}
		if !ok {
			c.JSON(common.HTTP_OK, common.Error("invalid expr result:"+Output.Expr))
			return
		}
	}

	OutputsBytes, _ := json.Marshal(AlarmRule.ExprDefine)
	Model := alarmcenter.MAlarmRule{
		UUID:        utils.AlarmRuleUuid(),
		Name:        AlarmRule.Name,
		ExprDefine:  string(OutputsBytes),
		Interval:    AlarmRule.Interval,
		Threshold:   AlarmRule.Threshold,
		HandleId:    AlarmRule.HandleId,
		Description: AlarmRule.Description,
	}
	if err := service.InsertAlarmRule(&Model); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ExprDefines := []alarmcenter.ExprDefine{}
	for _, ExprDefine := range AlarmRule.ExprDefine {
		ExprDefines = append(ExprDefines, alarmcenter.ExprDefine{
			Expr:      ExprDefine.Expr,
			EventType: ExprDefine.EventType,
		})
	}
	errLoadExpr := alarmcenter.LoadAlarmRule(Model.UUID, alarmcenter.AlarmRule{
		Interval:    time.Duration(AlarmRule.Interval) * time.Second,
		Threshold:   AlarmRule.Threshold,
		HandleId:    AlarmRule.HandleId,
		ExprDefines: ExprDefines,
	})
	if errLoadExpr != nil {
		c.JSON(common.HTTP_OK, common.Error400(errLoadExpr))
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
	for _, Output := range AlarmRule.ExprDefine {
		ok, errVerifyExpr := alarmcenter.VerifyExpr(Output.Expr)
		if errVerifyExpr != nil {
			c.JSON(common.HTTP_OK, common.Error400(errVerifyExpr))
			return
		}
		if !ok {
			c.JSON(common.HTTP_OK, common.Error("invalid expr result:"+Output.Expr))
			return
		}
	}

	OutputsBytes, _ := json.Marshal(AlarmRule.ExprDefine)

	Model := alarmcenter.MAlarmRule{
		UUID:        AlarmRule.UUID,
		Name:        AlarmRule.Name,
		ExprDefine:  string(OutputsBytes),
		Interval:    AlarmRule.Interval,
		Threshold:   AlarmRule.Threshold,
		HandleId:    AlarmRule.HandleId,
		Description: AlarmRule.Description,
	}
	if err := service.UpdateAlarmRule(&Model); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

// AlarmRuleTest
func AlarmRuleTest(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		Expr string         `json:"expr"`
		Data map[string]any `json:"data"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err.Error()))
		return
	}
	ok, err := alarmcenter.TestRunExpr(form.Expr, form.Data)
	if err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err.Error()))
		return
	}
	if !ok {
		c.JSON(common.HTTP_OK, common.OkWithData("invalid expr result:"+form.Expr))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData("SUCCESS"))
}
