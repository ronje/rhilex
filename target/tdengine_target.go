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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// http://<fqdn>:<port>/rest/sql/[db_name]
// fqnd: 集群中的任一台主机 FQDN 或 IP 地址
// port: 配置文件中 httpPort 配置项，缺省为 6041
// db_name: 可选参数，指定本次所执行的 SQL 语句的默认数据库库名
// curl -u root:taosdata -d 'show databases;' 106.15.225.172:6041/rest/sql
type TDEngineConfig struct {
	Fqdn             string `json:"fqdn" validate:"required" title:"地址"`     // 服务地址
	Port             int    `json:"port" validate:"required" title:"端口"`     // 服务端口
	Username         string `json:"username" validate:"required" title:"用户"` // 用户
	Password         string `json:"password" validate:"required" title:"密码"` // 密码
	DbName           string `json:"dbName" validate:"required" title:"数据库名"` // 数据库名
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}

/*
*
* TDengine 的资源输出支持,当前暂时支持HTTP接口的形式，后续逐步会增加UDP、TCP模式
*
 */

type tdEngineTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig TDEngineConfig
	status     typex.SourceState
}
type tdHttpResult struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Desc   string `json:"desc"`
}

func NewTdEngineTarget(e typex.Rhilex) typex.XTarget {
	td := tdEngineTarget{
		client: http.Client{Timeout: 2000 * time.Millisecond},
		mainConfig: TDEngineConfig{
			CacheOfflineData: new(bool),
		},
	}
	td.RuleEngine = e
	td.status = typex.SOURCE_DOWN
	return &td

}
func (td *tdEngineTarget) test() bool {
	if err := execQuery(td.client,
		td.mainConfig.Username,
		td.mainConfig.Password,
		"SELECT CLIENT_VERSION();",
		td.url()); err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	return true
}
func (td *tdEngineTarget) url() string {
	return fmt.Sprintf("http://%s:%v/rest/sql/%s",
		td.mainConfig.Fqdn, td.mainConfig.Port, td.mainConfig.DbName)
}

//
// 注册InEndID到资源
//

func (td *tdEngineTarget) Init(outEndId string, configMap map[string]interface{}) error {
	td.PointId = outEndId
	lostcache.CreateLostDataTable(outEndId)
	if err := utils.BindSourceConfig(configMap, &td.mainConfig); err != nil {
		return err
	}
	if td.test() {
		return nil
	}
	return errors.New("tdengine connect error")
}

// 启动资源
func (td *tdEngineTarget) Start(cctx typex.CCTX) error {
	td.Ctx = cctx.Ctx
	td.CancelCTX = cctx.CancelCTX
	//
	td.status = typex.SOURCE_UP
	// 补发数据
	if *td.mainConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(td.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				td.To(data.Data)
				{
					lostcache.DeleteLostCacheData(td.PointId, data.ID)
				}
			}
		}
	}

	return nil
}

// 获取资源状态
func (td *tdEngineTarget) Status() typex.SourceState {
	if td.test() {
		return typex.SOURCE_UP
	}
	return typex.SOURCE_DOWN
}

// 获取资源绑定的的详情
func (td *tdEngineTarget) Details() *typex.OutEnd {
	return td.RuleEngine.GetOutEnd(td.PointId)

}

// 停止资源, 用来释放资源
func (td *tdEngineTarget) Stop() {
	td.status = typex.SOURCE_DOWN
	if td.CancelCTX != nil {
		td.CancelCTX()
	}
}

func post(client http.Client,
	username string,
	password string,
	sql string,
	url string) (string, error) {
	body := strings.NewReader(sql)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "text/plain")
	request.SetBasicAuth(username, password)
	response, err2 := client.Do(request)
	if err2 != nil {
		return "", err2
	}
	if response.StatusCode != 200 {
		bytes0, err3 := io.ReadAll(response.Body)
		if err3 != nil {
			return "", err3
		}
		return "", fmt.Errorf("Error:%v", string(bytes0))
	}
	bytes1, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		return "", err3
	}
	return string(bytes1), nil
}

/*
*
* 执行TdEngine的查询
*
 */
func execQuery(client http.Client, username string, password string, sql string, url string) error {
	var r tdHttpResult
	// {"status":"error","code":534,"desc":"Syntax error in SQL"}
	glogger.GLogger.Debug("execQuery:", sql)

	body, err1 := post(client, username, password, sql, url)
	if err1 != nil {
		return err1
	}
	err2 := utils.TransformConfig([]byte(body), &r)
	if err2 != nil {
		return err2
	}
	if r.Status == "error" {
		return fmt.Errorf("code;%v, error:%s", r.Code, r.Desc)
	}
	return nil
}

/*
* SQL: INSERT INTO meter VALUES (NOW, %v, %v);
* 数据到达后写入Tdengine, 这里对数据有严格约束，必须是以,分割的字符串
* 比如: 10.22,220.12,123,......
*
 */
func (td *tdEngineTarget) To(data interface{}) (interface{}, error) {
	switch T := data.(type) {
	case string:
		{
			errQuery := execQuery(td.client, td.mainConfig.Username,
				td.mainConfig.Password, T, td.url())
			glogger.GLogger.Error(errQuery)
			if errQuery != nil {
				if *td.mainConfig.CacheOfflineData {
					lostcache.SaveLostCacheData(td.PointId, lostcache.CacheDataDto{
						TargetId: td.PointId,
						Data:     T,
					})
				}
			}
			return 0, errQuery
		}
	}
	return 0, nil
}
