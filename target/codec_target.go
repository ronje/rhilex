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

type codecTarget struct {
	typex.XStatus
	client        rhilexrpc.CodecClient
	rpcConnection *grpc.ClientConn
	mainConfig    GrpcConfig
	status        typex.SourceState
}

func NewCodecTarget(rx typex.Rhilex) typex.XTarget {
	ct := &codecTarget{}
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
func (ct *codecTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ct.PointId = outEndId
	//
	if err := utils.BindSourceConfig(configMap, &ct.mainConfig); err != nil {
		return err
	}

	return nil
}

// 启动资源
func (ct *codecTarget) Start(cctx typex.CCTX) error {
	ct.Ctx = cctx.Ctx
	ct.CancelCTX = cctx.CancelCTX
	//
	rpcConnection, err := grpc.NewClient(fmt.Sprintf("%s:%d", ct.mainConfig.Host, ct.mainConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	ct.rpcConnection = rpcConnection
	ct.client = rhilexrpc.NewCodecClient(rpcConnection)
	ct.status = typex.SOURCE_UP
	// 补发数据
	if CacheData, err1 := lostcache.GetLostCacheData(ct.PointId); err1 != nil {
		glogger.GLogger.Error(err1)
	} else {
		for _, data := range CacheData {
			_, errTo := ct.To(data.Data)
			if errTo == nil {
				lostcache.DeleteLostCacheData(data.ID)
			}
		}
	}

	return nil

}

// 获取资源状态
func (ct *codecTarget) Status() typex.SourceState {
	return ct.status

}

// 获取资源绑定的的详情
func (ct *codecTarget) Details() *typex.OutEnd {
	out := ct.RuleEngine.GetOutEnd(ct.PointId)
	return out

}

// 数据出口
func (ct *codecTarget) To(data interface{}) (interface{}, error) {
	dataRequest := &rhilexrpc.CodecRequest{
		Value: []byte(data.(string)),
	}
	var response *rhilexrpc.CodecResponse
	var err error
	if ct.mainConfig.Type == "DECODE" {
		response, err = ct.client.Decode(ct.Ctx, dataRequest)
	}
	if ct.mainConfig.Type == "ENCODE" {
		response, err = ct.client.Encode(ct.Ctx, dataRequest)
	}
	if err != nil {
		return nil, err
	}
	return response.GetData(), nil
}

// 停止资源, 用来释放资源
func (ct *codecTarget) Stop() {
	ct.status = typex.SOURCE_DOWN
	ct.CancelCTX()
	if ct.rpcConnection != nil {
		ct.rpcConnection.Close()
		ct.rpcConnection = nil
	}

}
