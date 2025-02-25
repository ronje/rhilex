package test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rhilex/component/apiserver"
	"github.com/hootrhino/rhilex/component/rhilexrpc"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/engine"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/plugin"
	"github.com/hootrhino/rhilex/plugin/demo_plugin"
	"github.com/hootrhino/rhilex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Modbus_LUA_Parse(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	engine := engine.NewRuleEngine(core.InitGlobalConfig("config/rhilex.ini"))
	engine.Start()
	hh := httpserver.NewHttpApiServer(engine)
	// HttpApiServer loaded default
	if err := plugin.LoadPlugin("plugin.http_server", hh); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}
	// Load a demo plugin
	if err := plugin.LoadPlugin("plugin.demo", demo_plugin.NewDemoPlugin()); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "rhilex Grpc InEnd", "rhilex Grpc InEnd", map[string]any{
		"port": 2581,
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}

	rule := typex.NewRule(engine,
		"uuid",
		"Just a test",
		"Just a test",
		[]string{grpcInend.UUID}[0],
		"",
		`function Success() print("[LUA Success Callback]=> OK") end`,
		`
		Actions = {
			function(args)
			    print(data)
				local json = require("json")
				local V6 = json.decode(data)
				local V7 = json.encode(rhilexlib:MB(">a:16 b:8 c:8", data, false))
				-- {"a":"0000000000000001","b":"00000000","c":"00000001"}
				print("[LUA Actions Callback, rhilex.MatchBinary] ==>", V7)
				return true, args
			end
		}`,
		`function Failed(error) print("[LUA Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer conn.Close()
	client := rhilexrpc.NewRhilexRpcClient(conn)

	for i := 0; i < 10; i++ {
		resp, err := client.Request(context.Background(), &rhilexrpc.RpcRequest{
			Value: fmt.Sprintf(`{"co2":10,"hum":30,"lex":22,"temp":100,"idx":%d}`, i),
		})
		if err != nil {
			glogger.GLogger.Errorf("grpc.Dial err: %v", err)
		}
		glogger.GLogger.Infof("rhilex Rpc Call Result ====>>: %v", resp.GetMessage())
	}
	time.Sleep(1 * time.Second)
	engine.Stop()
}
