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
	"fmt"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* GRPC
*
 */
type GrpcConfig struct {
	Host             string `json:"host" validate:"required" title:"地址"`
	Port             int    `json:"port" validate:"required" title:"端口"`
	Type             string `json:"type" title:"类型"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}

type RhilexRpcTarget struct {
	typex.XStatus
	client        rhilexrpc.RhilexRpcClient
	rpcConnection *grpc.ClientConn
	mainConfig    GrpcConfig
	status        typex.SourceState
}

func NewRhilexRpcTarget(rx typex.Rhilex) typex.XTarget {
	ct := &RhilexRpcTarget{}
	ct.mainConfig = GrpcConfig{
		Host:             "127.0.0.1",
		Port:             2581,
		CacheOfflineData: new(bool),
	}
	ct.RuleEngine = rx
	ct.status = typex.SOURCE_DOWN
	return ct
}

// 用来初始化传递资源配置
func (ct *RhilexRpcTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ct.PointId = outEndId
	lostcache.CreateLostDataTable(outEndId)
	//
	if err := utils.BindSourceConfig(configMap, &ct.mainConfig); err != nil {
		return err
	}
	return nil
}

// 启动资源
func (ct *RhilexRpcTarget) Start(cctx typex.CCTX) error {
	ct.Ctx = cctx.Ctx
	ct.CancelCTX = cctx.CancelCTX
	//
	rpcConnection, err := grpc.NewClient(fmt.Sprintf("%s:%d", ct.mainConfig.Host, ct.mainConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	ct.rpcConnection = rpcConnection
	ct.client = rhilexrpc.NewRhilexRpcClient(rpcConnection)
	ct.status = typex.SOURCE_UP
	// 补发数据
	if *ct.mainConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(ct.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				ct.To(data.Data)
				{
					lostcache.DeleteLostCacheData(ct.PointId, data.ID)
				}
			}
		}
	}

	return nil

}

// 获取资源状态
func (ct *RhilexRpcTarget) Status() typex.SourceState {
	return ct.status

}

// 获取资源绑定的的详情
func (ct *RhilexRpcTarget) Details() *typex.OutEnd {
	out := ct.RuleEngine.GetOutEnd(ct.PointId)
	return out

}

// 数据出口
func (ct *RhilexRpcTarget) To(data interface{}) (interface{}, error) {
	switch T := data.(type) {
	case string:
		dataRequest := &rhilexrpc.RpcRequest{
			Value: (T),
		}
		var err error
		_, err = ct.client.Request(ct.Ctx, dataRequest)

		if err != nil {
			if *ct.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ct.PointId, lostcache.CacheDataDto{
					TargetId: ct.PointId,
					Data:     data.(string),
				})
			}
			return 0, err
		}
		return 0, err
	}
	return 0, nil
}

// 停止资源, 用来释放资源
func (ct *RhilexRpcTarget) Stop() {
	ct.status = typex.SOURCE_DOWN
	if ct.CancelCTX != nil {
		ct.CancelCTX()
	}
	if ct.rpcConnection != nil {
		ct.rpcConnection.Close()
		ct.rpcConnection = nil
	}

}
