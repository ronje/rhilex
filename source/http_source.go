package source

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	"github.com/gin-gonic/gin"
	core "github.com/hootrhino/rhilex/config"
)

// httpInEndSource 表示HTTP输入端点源
type httpInEndSource struct {
	typex.XStatus
	engine     *gin.Engine
	mainConfig resconfig.HostConfig
	status     typex.SourceState
}

// NewHttpInEndSource 创建一个新的HTTP输入端点源
func NewHttpInEndSource(e typex.Rhilex) typex.XSource {
	h := httpInEndSource{}
	// 根据全局配置设置Gin的运行模式
	if core.GlobalConfig.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	h.engine = gin.New()
	h.RuleEngine = e
	return &h
}

// Init 初始化HTTP输入端点源
func (hh *httpInEndSource) Init(inEndId string, configMap map[string]any) error {
	hh.PointId = inEndId
	// 绑定配置
	if err := utils.BindSourceConfig(configMap, &hh.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind source config: %v", err)
		return err
	}
	// 验证配置
	if err := hh.validateConfig(); err != nil {
		glogger.GLogger.Errorf("Invalid config: %v", err)
		return err
	}
	return nil
}

// validateConfig 验证配置的有效性
func (hh *httpInEndSource) validateConfig() error {
	if hh.mainConfig.Port <= 0 || hh.mainConfig.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", hh.mainConfig.Port)
	}
	return nil
}

// Start 启动HTTP输入端点源
func (hh *httpInEndSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	// 注册POST路由
	hh.registerRoutes()

	// 启动HTTP服务器
	go func(ctx context.Context) {
		serverAddr := fmt.Sprintf(":%v", hh.mainConfig.Port)
		glogger.GLogger.Infof("Starting HTTP server on %s", serverAddr)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: hh.engine,
		}
		go func() {
			<-ctx.Done()
			glogger.GLogger.Info("Shutting down HTTP server...")
			if err := srv.Shutdown(context.Background()); err != nil {
				glogger.GLogger.Errorf("HTTP server shutdown error: %v", err)
			}
		}()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glogger.GLogger.Errorf("HTTP server error: %v", err)
		}
	}(hh.Ctx)

	hh.status = typex.SOURCE_UP
	glogger.GLogger.Infof("HTTP source started on [0.0.0.0]:%v", hh.mainConfig.Port)

	return nil
}

// registerRoutes 注册HTTP路由
func (hh *httpInEndSource) registerRoutes() {
	hh.engine.POST("/in", func(c *gin.Context) {
		// 定义请求体结构体
		type Form struct {
			Data string `json:"data"`
		}
		var inForm Form
		// 绑定JSON数据
		if err := c.BindJSON(&inForm); err != nil {
			glogger.GLogger.Errorf("Failed to bind JSON data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		// 调用RuleEngine处理数据
		hh.RuleEngine.WorkInEnd(hh.RuleEngine.GetInEnd(hh.PointId), inForm.Data)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"code":    http.StatusOK,
		})
	})
}

// Stop 停止HTTP输入端点源
func (hh *httpInEndSource) Stop() {
	hh.status = typex.SOURCE_DOWN
	if hh.CancelCTX != nil {
		hh.CancelCTX()
	}
}

// Status 获取HTTP输入端点源的状态
func (hh *httpInEndSource) Status() typex.SourceState {
	return hh.status
}

// Details 获取HTTP输入端点源的详细信息
func (hh *httpInEndSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}
