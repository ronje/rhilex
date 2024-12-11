package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/component/alarmcenter"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"gorm.io/gorm"
)

func LoadAlarmLogRoute() {
	api := server.RouteGroup(server.ContextUrl("/alarm_log"))
	api.GET(("/list"), server.AddRoute(AlarmLogList))
	api.PUT(("/update"), server.AddRoute(UpdateAlarmLog))
	api.DELETE(("/del"), server.AddRoute(DeleteAlarmLog))
	api.DELETE(("/clear"), server.AddRoute(ClearAlarmLog))
	api.GET(("/detail"), server.AddRoute(AlarmLogDetail))
}

// 告警日志
type AlarmLogVo struct {
	UUID      string `json:"uuid"`
	RuleId    string `json:"ruleId"`
	Source    string `json:"source"`
	EventType string `json:"eventType"`
	Ts        uint64 `json:"ts"`
	Summary   string `json:"summary"`
	Info      string `json:"info"`
}

/*
*
* AlarmLog
*
 */
func AlarmLogList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count := int64(0)
	var AlarmLogs = []alarmcenter.MAlarmLog{}
	ruleId, _ := c.GetQuery("ruleId")
	if ruleId != "" {
		ruleId, _ := c.GetQuery("ruleId")
		count, AlarmLogs = service.PageAlarmLogByRuleId(ruleId, pager.Current, pager.Size)
	} else {
		count, AlarmLogs = service.PageAlarmLog(pager.Current, pager.Size)
	}

	AlarmLogVos := []AlarmLogVo{}
	for _, AlarmLog := range AlarmLogs {
		AlarmLogVos = append(AlarmLogVos, AlarmLogVo{
			UUID:      AlarmLog.UUID,
			RuleId:    AlarmLog.RuleId,
			Source:    AlarmLog.Source,
			EventType: AlarmLog.EventType,
			Ts:        AlarmLog.Ts,
			Summary:   AlarmLog.Summary,
			Info:      AlarmLog.Info,
		})
	}
	Result := service.WrapPageResult(*pager, AlarmLogVos, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 详情
func AlarmLogDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	AlarmLog, err1 := service.GetMAlarmLogWithUUID(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err1))
		return
	}
	web_data := AlarmLogVo{
		UUID:      AlarmLog.UUID,
		RuleId:    AlarmLog.RuleId,
		Source:    AlarmLog.Source,
		EventType: AlarmLog.EventType,
		Ts:        AlarmLog.Ts,
		Summary:   AlarmLog.Summary,
		Info:      AlarmLog.Info,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(web_data))
}

/*
*
* 删除
*
 */
func DeleteAlarmLog(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	err := service.DeleteAlarmLog(uuid)
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

func CreateAlarmLog(c *gin.Context, ruleEngine typex.Rhilex) {
	AlarmLog := AlarmLogVo{}
	if err := c.ShouldBindJSON(&AlarmLog); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Model := alarmcenter.MAlarmLog{
		UUID:      utils.AlarmLogUuid(),
		RuleId:    AlarmLog.RuleId,
		Source:    AlarmLog.Source,
		EventType: AlarmLog.EventType,
		Ts:        AlarmLog.Ts,
		Summary:   AlarmLog.Summary,
		Info:      AlarmLog.Info,
	}
	if err := service.InsertAlarmLog(&Model); err != nil {
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

func UpdateAlarmLog(c *gin.Context, ruleEngine typex.Rhilex) {
	AlarmLog := AlarmLogVo{}
	if err := c.ShouldBindJSON(&AlarmLog); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Model := alarmcenter.MAlarmLog{
		UUID:      AlarmLog.UUID,
		RuleId:    AlarmLog.RuleId,
		Source:    AlarmLog.Source,
		EventType: AlarmLog.EventType,
		Ts:        AlarmLog.Ts,
		Summary:   AlarmLog.Summary,
		Info:      AlarmLog.Info,
	}
	if err := service.UpdateAlarmLog(&Model); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
func ClearAlarmLog(c *gin.Context, ruleEngine typex.Rhilex) {
	err := alarmcenter.AlarmDb().Transaction(func(tx *gorm.DB) error {
		tx.Exec("drop table m_alarm_logs if exists;")
		if tx.Error != nil {
			return tx.Error
		}
		alarmcenter.InitAlarmDbModel(tx)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	})
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
