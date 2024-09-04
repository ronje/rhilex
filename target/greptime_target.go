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
}
type GrepTimeDbTarget struct {
	typex.XStatus
	client     *greptime.Client
	table      *table.Table
	mainConfig GrepTimeDbTargetConfig
	status     typex.SourceState
}

func NewGrepTimeDbTarget(e typex.Rhilex) typex.XTarget {
	grep := new(GrepTimeDbTarget)
	grep.RuleEngine = e
	grep.mainConfig = GrepTimeDbTargetConfig{
		GwSn:     "rhilex",
		Host:     "127.0.0.1",
		Port:     4001,
		Username: "rhilex",
		Password: "rhilex",
		DataBase: "rhilex",
		Table:    "rhilex",
	}
	grep.status = typex.SOURCE_DOWN
	return grep
}

func (grep *GrepTimeDbTarget) Init(outEndId string, configMap map[string]interface{}) error {
	grep.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &grep.mainConfig); err != nil {
		return err
	}
	Table, err := table.New(grep.mainConfig.Table)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	grep.table = Table
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

	switch ST := data.(type) {
	case string:
		Map := map[string]interface{}{}
		errUnmarshal := json.Unmarshal([]byte(ST), &Map)
		if errUnmarshal != nil {
			glogger.GLogger.Error(errUnmarshal)
			return 0, errUnmarshal
		}
		grep.table.AddTimestampColumn("ts", types.TIMESTAMP_MILLISECOND)
		grep.table.AddTagColumn("device_name", types.STRING)
		values := []interface{}{time.Now().UnixMilli(), grep.mainConfig.GwSn}
		for k, v := range Map {
			switch VT := v.(type) {
			case bool:
				grep.table.AddTagColumn(k, types.BOOL)
				values = append(values, VT)
			case int32:
				grep.table.AddTagColumn(k, types.INT32)
				values = append(values, VT)
			case int64:
				grep.table.AddTagColumn(k, types.INT64)
				values = append(values, VT)
			case float32:
				grep.table.AddTagColumn(k, types.FLOAT32)
				values = append(values, VT)
			case float64:
				grep.table.AddTagColumn(k, types.FLOAT64)
				values = append(values, VT)
			case string:
				grep.table.AddTagColumn(k, types.STRING)
				values = append(values, VT)
			default:
				grep.table.AddTagColumn(k, types.STRING)
				values = append(values, VT)
			}
		}
		grep.table.AddRow(values...)
		glogger.GLogger.Debug("grep.client.Write: ", values)
		Response, errWrite := grep.client.Write(grep.Ctx, grep.table)
		if errWrite != nil {
			glogger.GLogger.Error(errWrite)
			return 0, errWrite
		}
		return Response.String(), nil
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
