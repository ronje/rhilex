package source

import (
	"context"
	"fmt"
	"net"

	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"google.golang.org/grpc"
)

const (
	defaultTransport = "tcp"
)

// GrpcConfig 定义gRPC服务器的配置
type GrpcConfig struct {
	Host             string `json:"host" validate:"required" title:"地址"`
	Port             int    `json:"port" validate:"required" title:"端口"`
	Type             string `json:"type" title:"类型"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}

// RhilexRpcServer 实现了RhilexRpc服务
type RhilexRpcServer struct {
	grpcInEndSource *grpcInEndSource
	rhilexrpc.UnimplementedRhilexRpcServer
}

// grpcInEndSource 表示gRPC输入端点源
type grpcInEndSource struct {
	typex.XStatus
	rhilexServer *RhilexRpcServer
	rpcServer    *grpc.Server
	mainConfig   GrpcConfig
	status       typex.SourceState
}

// NewGrpcInEndSource 创建一个新的gRPC输入端点源
func NewGrpcInEndSource(e typex.Rhilex) typex.XSource {
	g := grpcInEndSource{
		mainConfig: GrpcConfig{
			Host: "127.0.0.1",
			Port: 2583,
		},
	}
	g.RuleEngine = e
	return &g
}

// Init 初始化gRPC输入端点源
func (g *grpcInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	g.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &g.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind source config: %v", err)
		return err
	}
	if err := g.validateConfig(); err != nil {
		glogger.GLogger.Errorf("Invalid gRPC config: %v", err)
		return err
	}
	return nil
}

// validateConfig 验证gRPC配置的有效性
func (g *grpcInEndSource) validateConfig() error {
	if g.mainConfig.Port <= 0 || g.mainConfig.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", g.mainConfig.Port)
	}
	if g.mainConfig.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	return nil
}

// Start 启动gRPC输入端点源
func (g *grpcInEndSource) Start(cctx typex.CCTX) error {
	g.Ctx = cctx.Ctx
	g.CancelCTX = cctx.CancelCTX

	listener, err := net.Listen(defaultTransport, fmt.Sprintf(":%d", g.mainConfig.Port))
	if err != nil {
		glogger.GLogger.Errorf("Failed to listen on port %d: %v", g.mainConfig.Port, err)
		return err
	}

	g.rpcServer = grpc.NewServer()
	g.rhilexServer = &RhilexRpcServer{grpcInEndSource: g}
	rhilexrpc.RegisterRhilexRpcServer(g.rpcServer, g.rhilexServer)

	go func() {
		glogger.GLogger.Infof("RhilexRpc source started on %s", listener.Addr().String())
		if err := g.rpcServer.Serve(listener); err != nil {
			glogger.GLogger.Errorf("Failed to serve gRPC server: %v", err)
		}
	}()

	g.status = typex.SOURCE_UP
	return nil
}

// Stop 停止gRPC输入端点源
func (g *grpcInEndSource) Stop() {
	g.status = typex.SOURCE_DOWN
	if g.CancelCTX != nil {
		g.CancelCTX()
	}
	if g.rpcServer != nil {
		g.rpcServer.GracefulStop()
		g.rpcServer = nil
	}
}

// Status 获取gRPC输入端点源的状态
func (g *grpcInEndSource) Status() typex.SourceState {
	return g.status
}

// Details 获取gRPC输入端点源的详细信息
func (g *grpcInEndSource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}

// Request 处理gRPC请求
func (r *RhilexRpcServer) Request(ctx context.Context, in *rhilexrpc.RpcRequest) (*rhilexrpc.RpcResponse, error) {
	ok, err := r.grpcInEndSource.RuleEngine.WorkInEnd(
		r.grpcInEndSource.RuleEngine.GetInEnd(r.grpcInEndSource.PointId),
		string(in.Value),
	)
	if err != nil {
		glogger.GLogger.Errorf("Failed to process gRPC request: %v", err)
		return &rhilexrpc.RpcResponse{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	if ok {
		return &rhilexrpc.RpcResponse{
			Code:    0,
			Message: "OK",
		}, nil
	}
	return &rhilexrpc.RpcResponse{
		Code:    1,
		Message: "Request processing failed",
	}, nil
}
