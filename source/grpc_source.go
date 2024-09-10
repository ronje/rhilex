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

type GrpcConfig struct {
	Host             string `json:"host" validate:"required" title:"地址"`
	Port             int    `json:"port" validate:"required" title:"端口"`
	Type             string `json:"type" title:"类型"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}
type RhilexRpcServer struct {
	grpcInEndSource *grpcInEndSource
	rhilexrpc.UnimplementedRhilexRpcServer
}

// Source interface
type grpcInEndSource struct {
	typex.XStatus
	rhilexServer *RhilexRpcServer
	rpcServer    *grpc.Server
	mainConfig   GrpcConfig
	status       typex.SourceState
}

func NewGrpcInEndSource(e typex.Rhilex) typex.XSource {
	g := grpcInEndSource{}
	g.RuleEngine = e
	g.mainConfig = GrpcConfig{}
	return &g
}

/*
*
* Init
*
 */
func (g *grpcInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	g.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &g.mainConfig); err != nil {
		return err
	}
	return nil
}

func (g *grpcInEndSource) Start(cctx typex.CCTX) error {
	g.Ctx = cctx.Ctx
	g.CancelCTX = cctx.CancelCTX

	listener, err := net.Listen(defaultTransport, fmt.Sprintf(":%d", g.mainConfig.Port))
	if err != nil {
		return err
	}
	// Important !!!
	g.rpcServer = grpc.NewServer()
	g.rhilexServer = new(RhilexRpcServer)
	g.rhilexServer.grpcInEndSource = g
	//
	rhilexrpc.RegisterRhilexRpcServer(g.rpcServer, g.rhilexServer)

	go func(c context.Context) {
		glogger.GLogger.Info("RhilexRpc source started on", listener.Addr())
		g.rpcServer.Serve(listener)
	}(g.Ctx)
	g.status = typex.SOURCE_UP
	return nil
}

func (g *grpcInEndSource) Stop() {
	g.status = typex.SOURCE_DOWN
	if g.CancelCTX != nil {
		g.CancelCTX()
	}
	if g.rpcServer != nil {
		g.rpcServer.Stop()
		g.rpcServer = nil
	}

}

func (g *grpcInEndSource) Status() typex.SourceState {
	return g.status
}

func (g *grpcInEndSource) Test(inEndId string) bool {
	return true
}

func (g *grpcInEndSource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}

func (r *RhilexRpcServer) Work(ctx context.Context, in *rhilexrpc.Data) (*rhilexrpc.Response, error) {
	ok, err := r.grpcInEndSource.RuleEngine.WorkInEnd(
		r.grpcInEndSource.RuleEngine.GetInEnd(r.grpcInEndSource.PointId),
		in.Value,
	)
	if ok {
		return &rhilexrpc.Response{
			Code:    0,
			Message: "OK",
		}, nil
	} else {
		return &rhilexrpc.Response{
			Code:    1,
			Message: err.Error(),
		}, err
	}

}

// 来自外面的数据
func (*grpcInEndSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*grpcInEndSource) UpStream([]byte) (int, error) {
	return 0, nil
}
