package source

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	coap "github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
)

// coAPInEndSource 表示一个基于CoAP协议的输入端点源
type coAPInEndSource struct {
	typex.XStatus
	router     *mux.Router
	mainConfig resconfig.HostConfig
	status     typex.SourceState
}

// NewCoAPInEndSource 创建一个新的CoAP输入端点源
func NewCoAPInEndSource(e typex.Rhilex) typex.XSource {
	c := coAPInEndSource{
		router: mux.NewRouter(),
		// 初始化时可以考虑设置一些默认配置
		mainConfig: resconfig.HostConfig{
			Host: "127.0.0.1",
			Port: 2584,
		},
	}
	c.RuleEngine = e
	return &c
}

// Init 初始化CoAP输入端点源
func (cc *coAPInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	cc.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &cc.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind source config: %v", err)
		return err
	}
	// 可以在这里添加配置验证逻辑
	if cc.mainConfig.Port <= 0 || cc.mainConfig.Port > 65535 {
		glogger.GLogger.Errorf("Invalid port number: %d", cc.mainConfig.Port)
		return fmt.Errorf("invalid port number: %d", cc.mainConfig.Port)
	}
	return nil
}

// InData 表示输入的数据结构
type InData struct {
	Ts      uint64 `json:"ts"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Start 启动CoAP输入端点源
func (cc *coAPInEndSource) Start(cctx typex.CCTX) error {
	cc.Ctx = cctx.Ctx
	cc.CancelCTX = cctx.CancelCTX
	port := fmt.Sprintf(":%v", cc.mainConfig.Port)

	// 定义处理函数，提高代码可读性
	handleCoAPRequest := func(w mux.ResponseWriter, msg *mux.Message) {
		Body, err := msg.ReadBody()
		if err != nil {
			glogger.GLogger.Errorf("Failed to read CoAP message body: %v", err)
			return
		}

		payload := InData{
			Ts:      uint64(time.Now().UnixMilli()),
			Type:    "POST",
			Payload: string(Body),
		}

		bites, err := json.Marshal(payload)
		if err != nil {
			glogger.GLogger.Errorf("Failed to marshal CoAP payload: %v", err)
			return
		}

		glogger.GLogger.Debug(string(bites))
		cc.RuleEngine.WorkInEnd(cc.RuleEngine.GetInEnd(cc.PointId), string(bites))

		if err := w.SetResponse(codes.Content, message.AppOctets, bytes.NewReader([]byte{200})); err != nil {
			glogger.GLogger.Errorf("Cannot set CoAP response: %v", err)
		}
	}

	cc.router.Handle("/", mux.HandlerFunc(handleCoAPRequest))

	// 启动CoAP服务器
	go func(ctx context.Context) {
		glogger.GLogger.Infof("Starting CoAP server on port %s", port)
		if err := coap.ListenAndServe("udp", port, cc.router); err != nil {
			glogger.GLogger.Errorf("Failed to start CoAP server: %v", err)
		}
	}(cc.Ctx)

	cc.status = typex.SOURCE_UP
	glogger.GLogger.Infof("Coap source started on [udp]%s", port)
	return nil
}

// Stop 停止CoAP输入端点源
func (cc *coAPInEndSource) Stop() {
	cc.status = typex.SOURCE_DOWN
	if cc.CancelCTX != nil {
		cc.CancelCTX()
	}
	// 可以考虑添加一些资源清理逻辑，比如关闭CoAP服务器
}

// Status 获取CoAP输入端点源的状态
func (cc *coAPInEndSource) Status() typex.SourceState {
	return cc.status
}

// Details 获取CoAP输入端点源的详细信息
func (cc *coAPInEndSource) Details() *typex.InEnd {
	return cc.RuleEngine.GetInEnd(cc.PointId)
}
