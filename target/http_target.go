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
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type HTTPTargetConfig struct {
	Url              string            `json:"url" validate:"required" title:"URL"`
	Headers          map[string]string `json:"headers" validate:"required" title:"HTTP Headers"`
	AllowPing        *bool             `json:"allowPing"`
	PingPacket       string            `json:"pingPacket"`
	Timeout          int               `json:"timeout"`
	CacheOfflineData *bool             `json:"cacheOfflineData" title:"离线缓存"`
}

type HTTPTargetMainConfig struct {
	HTTPTargetConfig HTTPTargetConfig `json:"commonConfig" validate:"required"`
}
type HTTPTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig HTTPTargetMainConfig
	status     typex.SourceState
}

func NewHTTPTarget(e typex.Rhilex) typex.XTarget {
	ht := new(HTTPTarget)
	ht.RuleEngine = e
	ht.mainConfig = HTTPTargetMainConfig{
		HTTPTargetConfig: HTTPTargetConfig{
			Url:              "http://127.0.0.1",
			PingPacket:       "rhilex",
			Timeout:          3000,
			AllowPing:        new(bool),
			Headers:          map[string]string{},
			CacheOfflineData: new(bool),
		},
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *HTTPTarget) Init(outEndId string, configMap map[string]any) error {
	ht.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}

	return nil

}
func (ht *HTTPTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	ht.client = http.Client{}
	ht.status = typex.SOURCE_UP
	if *ht.mainConfig.HTTPTargetConfig.AllowPing {
		go func(ht *HTTPTarget) {
			for {
				select {
				case <-ht.Ctx.Done():
					return
				default:
				}
				_, err := utils.Post(ht.client, ht.mainConfig.HTTPTargetConfig.PingPacket,
					ht.mainConfig.HTTPTargetConfig.Url, ht.mainConfig.HTTPTargetConfig.Headers)
				if err != nil {
					glogger.GLogger.Error(err)
					ht.status = typex.SOURCE_DOWN
					continue
				}
				time.Sleep(time.Duration(ht.mainConfig.HTTPTargetConfig.Timeout) * time.Millisecond)
			}
		}(ht)
	}
	// 补发数据
	if *ht.mainConfig.HTTPTargetConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(ht.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				ht.To(data.Data)
				{
					lostcache.DeleteLostCacheData(ht.PointId, data.ID)
				}
			}
		}
	}
	glogger.GLogger.Info("HTTP Target started")
	return nil
}

func (ht *HTTPTarget) Status() typex.SourceState {
	if err := ht.prob(); err != nil {
		glogger.GLogger.Error(err)
		return typex.SOURCE_DOWN
	}
	return ht.status

}

func (ht *HTTPTarget) To(data any) (any, error) {
	switch T := data.(type) {
	case string:

		_, err := utils.Post(ht.client, T,
			ht.mainConfig.HTTPTargetConfig.Url, ht.mainConfig.HTTPTargetConfig.Headers)
		if err != nil {
			glogger.GLogger.Error(err)
			if *ht.mainConfig.HTTPTargetConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ht.PointId, lostcache.CacheDataDto{
					TargetId: ht.PointId,
					Data:     T,
				})
			}
			return nil, err
		}
	}
	return nil, fmt.Errorf("data type must string!")
}

func (ht *HTTPTarget) Stop() {
	ht.status = typex.SOURCE_DOWN
	if ht.CancelCTX != nil {
		ht.CancelCTX()
	}
}
func (ht *HTTPTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}
func (ht *HTTPTarget) prob() error {
	d := net.Dialer{
		Timeout: 3 * time.Second,
	}
	Url, err := url.Parse(ht.mainConfig.HTTPTargetConfig.Url)
	if err != nil {
		return err
	}
	conn, err := d.Dial("tcp", Url.Host)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
