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

package target

import (
	"encoding/json"
	"time"

	"github.com/GreptimeTeam/greptimedb-ingester-go/table"
	"github.com/GreptimeTeam/greptimedb-ingester-go/table/types"
	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	greptime "github.com/GreptimeTeam/greptimedb-ingester-go"
)

type GrepTimeDbTargetConfig struct {
	GwSn     string `json:"gwsn" validate:"required" title:"序列号"`      // 服务地址
	Host     string `json:"host" validate:"required" title:"地址"`       // 服务地址
	Port     int    `json:"port" validate:"required" title:"端口"`       // 服务端口
	Username string `json:"username" validate:"required" title:"用户"`   // 用户
	Password string `json:"password" validate:"required" title:"密码"`   // 密码
	DataBase string `json:"database" validate:"required" title:"数据库名"` // 数据库名
	Table    string `json:"table" validate:"required" title:"数据表"`     // 表名
	// 离线缓存
	CacheOfflineData *bool `json:"cacheOfflineData" title:"离线缓存"`
}
type GrepTimeDbTarget struct {
	typex.XStatus
	client     *greptime.Client
	mainConfig GrepTimeDbTargetConfig
	status     typex.SourceState
}

func NewGrepTimeDbTarget(e typex.Rhilex) typex.XTarget {
	grep := new(GrepTimeDbTarget)
	grep.RuleEngine = e
	grep.mainConfig = GrepTimeDbTargetConfig{
		GwSn:             "rhilex",
		Host:             "127.0.0.1",
		Port:             4001,
		Username:         "rhilex",
		Password:         "rhilex",
		DataBase:         "public",
		Table:            "rhilex",
		CacheOfflineData: new(bool),
	}
	grep.status = typex.SOURCE_DOWN
	return grep
}

func (grep *GrepTimeDbTarget) Init(outEndId string, configMap map[string]interface{}) error {
	grep.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &grep.mainConfig); err != nil {
		return err
	}
	return nil

}
func (grep *GrepTimeDbTarget) Start(cctx typex.CCTX) error {
	grep.Ctx = cctx.Ctx
	grep.CancelCTX = cctx.CancelCTX
	//
	cfg := greptime.NewConfig(grep.mainConfig.Host).WithPort(grep.mainConfig.Port).
		WithAuth(grep.mainConfig.Username, grep.mainConfig.Password).
		WithDatabase(grep.mainConfig.DataBase)
	cfg.WithKeepalive(time.Second*10, time.Second*5)
	client, err := greptime.NewClient(cfg)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	grep.client = client
	grep.status = typex.SOURCE_UP
	// 补发数据
	if *grep.mainConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(grep.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				_, errTo := grep.To(data.Data)
				if errTo == nil {
					lostcache.DeleteLostCacheData(data.ID)
				}
			}
		}

	}

	glogger.GLogger.Info("Template Target started")
	return nil
}

func (grep *GrepTimeDbTarget) Status() typex.SourceState {
	if grep.client != nil {
		return typex.SOURCE_UP
	}
	return grep.status

}

// To: data-Map
func (grep *GrepTimeDbTarget) To(data interface{}) (interface{}, error) {

	switch T := data.(type) {
	case string:
		Map := map[string]interface{}{}
		errUnmarshal := json.Unmarshal([]byte(T), &Map)
		if errUnmarshal != nil {
			glogger.GLogger.Error(errUnmarshal)
			return 0, errUnmarshal
		}
		Table, errNew := table.New(grep.mainConfig.Table)
		if errNew != nil {
			glogger.GLogger.Error(errNew)
			return 0, errNew
		}
		Table.AddTimestampColumn("ts", types.TIMESTAMP_MILLISECOND)
		Table.AddTagColumn("gateway_sn", types.STRING)
		values := []interface{}{time.Now().UnixMilli(), grep.mainConfig.GwSn}
		for k, v := range Map {
			switch VT := v.(type) {
			case bool:
				Table.AddFieldColumn(k, types.BOOL)
				values = append(values, VT)
			case int32:
				Table.AddFieldColumn(k, types.INT32)
				values = append(values, VT)
			case int64:
				Table.AddFieldColumn(k, types.INT64)
				values = append(values, VT)
			case float32:
				Table.AddFieldColumn(k, types.FLOAT32)
				values = append(values, VT)
			case float64:
				Table.AddFieldColumn(k, types.FLOAT64)
				values = append(values, VT)
			case string:
				Table.AddFieldColumn(k, types.STRING)
				values = append(values, VT)
			default:
				Table.AddFieldColumn(k, types.STRING)
				values = append(values, VT)
			}
		}
		Table.AddRow(values...)
		_, errWrite := grep.client.Write(grep.Ctx, Table)
		glogger.GLogger.Debug("grep.client.Write: ", values)
		if errWrite != nil {
			glogger.GLogger.Error(errWrite)
			if *grep.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(lostcache.CacheDataDto{
					TargetId: grep.PointId,
					Data:     T,
				})
			}
			return 0, errWrite
		}
		return 0, errWrite
	}
	return 0, nil
}

func (grep *GrepTimeDbTarget) Stop() {
	grep.status = typex.SOURCE_DOWN
	if grep.CancelCTX != nil {
		grep.CancelCTX()
	}
	if grep.client != nil {
		grep.client = nil
	}
}
func (grep *GrepTimeDbTarget) Details() *typex.OutEnd {
	return grep.RuleEngine.GetOutEnd(grep.PointId)
}
