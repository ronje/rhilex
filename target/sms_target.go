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

package target

import (
	"encoding/json"
	"fmt"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	gosms "github.com/pkg6/go-sms"
	"github.com/pkg6/go-sms/gateways/aliyun"
	"github.com/pkg6/go-sms/gateways/juhe"
)

/**
 *
 * 短信
 */
type SMSTargetConfig struct {
	Type string `json:"type" validate:"required"` // ALI_SMS|JUHE_SMS
	// juhe
	AppId  string `json:"app_id"`
	AppKey string `json:"app_key"`
	// aliyun
	AccessKeyId      string `json:"accessKeyId"`
	AccessKeySecret  string `json:"accessKeySecret"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}
type SMSTarget struct {
	typex.XStatus
	mainConfig SMSTargetConfig
	status     typex.SourceState
}

func NewSMSTarget(e typex.Rhilex) typex.XTarget {
	ht := new(SMSTarget)
	ht.RuleEngine = e
	ht.mainConfig = SMSTargetConfig{Type: "default"}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *SMSTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}
	return nil

}
func (ht *SMSTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	ht.status = typex.SOURCE_UP
	glogger.GLogger.Info("Template Target started")
	return nil
}

func (ht *SMSTarget) Status() typex.SourceState {
	return ht.status
}

type sms_template struct {
	To       int               `json:"to"`
	Content  string            `json:"content"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

func (O sms_template) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}
func (ht *SMSTarget) To(data interface{}) (interface{}, error) {
	sms_template := sms_template{}
	switch T := data.(type) {
	case string:
		if err := json.Unmarshal([]byte(T), &sms_template); err != nil {
			return nil, err
		}
		glogger.GLogger.Debug("Send SMS:", sms_template.String())
		var gateway gosms.IGateway
		if ht.mainConfig.Type == "ALI_SMS" {
			gateway = aliyun.GateWay(ht.mainConfig.AccessKeyId, ht.mainConfig.AccessKeySecret)
		}
		if ht.mainConfig.Type == "JUHE_SMS" {
			gateway = juhe.GateWay(ht.mainConfig.AppKey)
		}
		if gateway != nil {
			number := gosms.NoCodePhoneNumber(sms_template.To)
			message := gosms.MessageTemplate(sms_template.Template, sms_template.Data)
			result, err := gosms.Sender(number, message, gateway)
			if err != nil {
				glogger.GLogger.Error("gosms.Sender:", err, result.ClientResult)
				return nil, err
			}
			if ht.mainConfig.Type == "ALI_SMS" {
				if resp, ok := result.ClientResult.Response.(aliyun.Response); !ok {
					return 0, fmt.Errorf(resp.Message)
				}
			}
			if ht.mainConfig.Type == "JUHE_SMS" {
				if resp, ok := result.ClientResult.Response.(juhe.Response); !ok {
					return 0, fmt.Errorf(resp.Reason)
				}
			}
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("sms format error")
	}
}

func (ht *SMSTarget) Stop() {
	ht.status = typex.SOURCE_DOWN
	if ht.CancelCTX != nil {
		ht.CancelCTX()
	}
}
func (ht *SMSTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}
