// Copyright (C) 2024 Your Name
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

	"github.com/go-redis/redis/v8"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// RedisTargetConfig 用于存储RedisTarget的配置信息
type RedisTargetConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// RedisTarget 实现了将数据存储到Redis的目标组件
type RedisTarget struct {
	typex.XStatus
	mainConfig RedisTargetConfig
	status     typex.SourceState
	client     *redis.Client
}

// RedisData 用于封装要存储到Redis的数据和键
type RedisData struct {
	Key  string
	Data map[string]interface{}
}

// NewRedisTarget 创建一个新的RedisTarget实例
func NewRedisTarget(e typex.Rhilex) typex.XTarget {
	rt := new(RedisTarget)
	rt.RuleEngine = e
	rt.mainConfig = RedisTargetConfig{}
	rt.status = typex.SOURCE_DOWN
	return rt
}

// Init 初始化RedisTarget
func (rt *RedisTarget) Init(outEndId string, configMap map[string]interface{}) error {
	rt.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &rt.mainConfig); err != nil {
		return err
	}

	// 初始化Redis客户端
	rt.client = redis.NewClient(&redis.Options{
		Addr:     rt.mainConfig.Address,
		Password: rt.mainConfig.Password,
		DB:       rt.mainConfig.DB,
	})

	// 测试Redis连接
	_, err := rt.client.Ping(rt.client.Context()).Result()
	if err != nil {
		return err
	}

	return nil
}

// Start 启动RedisTarget
func (rt *RedisTarget) Start(cctx typex.CCTX) error {
	rt.Ctx = cctx.Ctx
	rt.CancelCTX = cctx.CancelCTX
	rt.status = typex.SOURCE_UP
	glogger.GLogger.Info("Redis Target started")
	return nil
}

// Status 获取RedisTarget的当前状态，返回Redis PING的结果
func (rt *RedisTarget) Status() typex.SourceState {
	// 发送PING命令
	pong, err := rt.client.Ping(rt.client.Context()).Result()
	if err != nil {
		rt.status = typex.SOURCE_DOWN
		glogger.GLogger.Error("Redis connection error:", err)
	} else if pong == "PONG" {
		rt.status = typex.SOURCE_UP
	}
	return rt.status
}

// To 将数据存储到Redis中，使用 HMSET 命令
func (rt *RedisTarget) To(data interface{}) (interface{}, error) {
	// 将 data 转换为 RedisData 结构体
	redisData, ok := data.(RedisData)
	if !ok {
		return nil, fmt.Errorf("input data is not of type RedisData")
	}

	// 将数据存储到Redis中，使用 HMSET 命令
	err := rt.client.HMSet(rt.client.Context(), redisData.Key, redisData.Data).Err()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Stop 停止RedisTarget
func (rt *RedisTarget) Stop() {
	rt.status = typex.SOURCE_DOWN
	if rt.CancelCTX != nil {
		rt.CancelCTX()
	}
	// 关闭Redis客户端
	if rt.client != nil {
		rt.client.Close()
	}
}

// Details 获取RedisTarget关联的输出端点的详细信息
func (rt *RedisTarget) Details() *typex.OutEnd {
	return rt.RuleEngine.GetOutEnd(rt.PointId)
}
