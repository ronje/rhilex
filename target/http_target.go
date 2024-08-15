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
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type HTTPTargetConfig struct {
	Url        string            `json:"url" validate:"required" title:"URL"`
	Headers    map[string]string `json:"headers" validate:"required" title:"HTTP Headers"`
	AllowPing  *bool             `json:"allowPing"`
	PingPacket string            `json:"pingPacket"`
	Timeout    int               `json:"timeout"`
}
type HTTPTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig HTTPTargetConfig
	status     typex.SourceState
}

func NewHTTPTarget(e typex.Rhilex) typex.XTarget {
	ht := new(HTTPTarget)
	ht.RuleEngine = e
	ht.mainConfig = HTTPTargetConfig{
		PingPacket: "rhilex",
		Timeout:    3000,
		AllowPing: func() *bool {
			b := true
			return &b
		}(),
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *HTTPTarget) Init(outEndId string, configMap map[string]interface{}) error {
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
	if *ht.mainConfig.AllowPing {
		go func(ht *HTTPTarget) {
			for {
				select {
				case <-ht.Ctx.Done():
					return
				default:
				}
				_, err := utils.Post(ht.client, ht.mainConfig.PingPacket,
					ht.mainConfig.Url, ht.mainConfig.Headers)
				if err != nil {
					glogger.GLogger.Error(err)
					ht.status = typex.SOURCE_DOWN
					continue
				}
				time.Sleep(time.Duration(ht.mainConfig.Timeout) * time.Millisecond)
			}
		}(ht)
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

type HTTPTargetOutputData struct {
	Label string `json:"label"`
	Body  string `json:"body"`
}

func (O HTTPTargetOutputData) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}
func (ht *HTTPTarget) To(data interface{}) (interface{}, error) {
	switch T := data.(type) {
	case string:
		outputData := HTTPTargetOutputData{
			Label: ht.mainConfig.PingPacket,
			Body:  T,
		}
		_, err := utils.Post(ht.client, outputData.String(),
			ht.mainConfig.Url, ht.mainConfig.Headers)
		if err != nil {
			glogger.GLogger.Error(err)
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
	Url, err := url.Parse(ht.mainConfig.Url)
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
