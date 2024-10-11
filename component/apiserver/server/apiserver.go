package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	response "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 拼接URL
*
 */
func ContextUrl(path string) string {
	return API_V1_ROOT + path
}

const API_V1_ROOT string = "/api/v1/"

// DefaultApiServer 全局API Server
var DefaultApiServer *RhilexApiServer

/*
*
* API Server
*
 */
type RhilexApiServer struct {
	ginEngine  *gin.Engine
	ruleEngine typex.Rhilex
	config     serverConfig
}
type serverConfig struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

/*
*
* 开启Server
*
 */
func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}
func errorHandler(c *gin.Context, info Info) {
	c.JSON(400, response.Error("Too many requests. Try again after 3s"))
}
func StartRhilexApiServer(ruleEngine typex.Rhilex, port int) {
	gin.SetMode(gin.ReleaseMode)
	// if core.GlobalConfig.AppDebugMode {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	server := RhilexApiServer{
		ginEngine:  gin.New(),
		ruleEngine: ruleEngine,
		config:     serverConfig{Port: port},
	}
	RateLimiter := RateLimiter(InMemoryStore(&InMemoryOptions{
		Rate: time.Second, Limit: 30}), &Options{
		ErrorHandler: errorHandler, KeyFunc: keyFunc,
	})
	server.ginEngine.Use(RateLimiter)
	staticFs := WWWRoot("")
	// Logo 静态资源路由
	server.ginEngine.Use(func(ctx *gin.Context) {
		if ctx.Request.RequestURI == "/logo.svg" {
			ctx.Header("Content-Type", "image/svg+xml")
			ctx.FileFromFS("logo.svg", staticFs)
			ctx.Writer.Flush()
			return
		}
		if ctx.Request.RequestURI == "/favicon.svg" {
			ctx.Header("Content-Type", "image/svg+xml")
			ctx.FileFromFS("favicon.svg", staticFs)
			ctx.Writer.Flush()
			return
		}
		ctx.Next()
	})
	server.ginEngine.Use(static.Serve("/", staticFs))
	server.ginEngine.Use(Authorize())
	server.ginEngine.Use(CheckLicense())
	server.ginEngine.Use(Cros())
	server.ginEngine.GET("/ws", glogger.WsLogger)
	server.ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		if core.GlobalConfig.AppDebugMode {
			debug.PrintStack()
			os.Exit(1)
			panic(err)
		}
		c.JSON(500, response.Error500(errors.New("http server crash, try to recovery")))
	}))
	// 解决浏览器刷新被重定向问题
	server.ginEngine.NoRoute(func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.JSON(404, response.Error("No such Route:"+c.Request.URL.Path))
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Write(indexHTML)
		c.Writer.Flush()
	})
	go func(ctx context.Context, port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			glogger.GLogger.Fatalf("Http Api Server listen error: %s", err)
		}
		defer listener.Close()
		if err := server.ginEngine.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("Http Api Server listen error: %s", err)
		}
	}(typex.GCTX, server.config.Port)
	glogger.GLogger.Infof("Http Api Server listen on: %s", fmt.Sprintf("0.0.0.0:%d", port))
	DefaultApiServer = &server
}

// 即将废弃
func (s *RhilexApiServer) AddRoute(f func(c *gin.Context,
	ruleEngine typex.Rhilex)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, s.ruleEngine)
	}
}

// New api after 0.6.4
func AddRoute(f func(c *gin.Context,
	ruleEngine typex.Rhilex)) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c, DefaultApiServer.ruleEngine)
	}
}

// AddRouteV2: Only for cron，It's redundant API， will refactor in feature
func AddRouteV2(f func(*gin.Context, typex.Rhilex) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		data, err := f(c, DefaultApiServer.ruleEngine)
		if err != nil {
			c.JSON(response.HTTP_OK, response.Error400(err))
		} else {
			c.JSON(response.HTTP_OK, response.OkWithData(data))
		}
	}
}

func (s *RhilexApiServer) GetGroup(name string) *gin.RouterGroup {
	return s.ginEngine.Group(name)
}

// new API
func RouteGroup(name string) *gin.RouterGroup {
	return DefaultApiServer.ginEngine.Group(name)
}
func (s *RhilexApiServer) Route() *gin.Engine {
	return s.ginEngine
}

/*
*
* 初始化网络配置,主要针对Linux，而且目前只支持了Ubuntu1804,后期分一个windows版本
*
 */
func (s *RhilexApiServer) InitializeGenericOSData() {
	initStaticModel()
}
func (s *RhilexApiServer) InitializeUnixData() {
	glogger.GLogger.Info("Initialize Unix Default Data")
}
func (s *RhilexApiServer) InitializeWindowsData() {
	glogger.GLogger.Info("Initialize Windows Default Data")
}

func (s *RhilexApiServer) InitializeProduct() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "RHILEXG1" {
		glogger.GLogger.Info("Initialize Rhilex Pi Default Data")
	}
}

/*
*
* 初始化一些静态数据
*
 */
func initStaticModel() {
	// 配置一个默认分组
	service.InitGenericGroup(&model.MGenericGroup{
		UUID:   "DROOT",
		Type:   "DEVICE",
		Name:   "DefaultGroup",
		Parent: "NULL",
	})
	service.InitGenericGroup(&model.MGenericGroup{
		UUID:   "ULTROOT",
		Type:   "USER_LUA_TEMPLATE",
		Name:   "DefaultGroup",
		Parent: "NULL",
	})
	// 初始化一个用户
	service.InitMUser(
		&model.MUser{
			Role:        "Admin",
			Username:    "rhilex",
			Password:    "25d55ad283aa400af464c76d713c07ad", // md5(12345678)
			Description: "Default RHILEX Admin User",
		},
	)
}
