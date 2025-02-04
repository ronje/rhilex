package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rhilex/component/apiserver"
	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hootrhino/rhilex/typex"
)

// go test -timeout 30s -run ^Test_DataToSemtechUdp github.com/hootrhino/rhilex/test -v -count=1

func Test_DataToSemtechUdp(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	if err := registry.DefaultPluginRegistry.LoadPlugin("plugin.http_server",
		httpserver.NewHttpApiServer(engine)); err != nil {
		t.Fatal("HttpServer load failed:", err)
	}

	//
	SEMTECH_UDP_FORWARDER := typex.NewOutEnd(typex.SEMTECH_UDP_FORWARDER,
		"SEMTECH_UDP_FORWARDER", "SEMTECH_UDP_FORWARDER", map[string]interface{}{
			"host":             "192.168.10.163",
			"port":             1700,
			"mac":              "a46a4de31a346180",
			"cacheOfflineData": true,
		},
	)
	SEMTECH_UDP_FORWARDER.UUID = "SEMTECH_UDP_FORWARDER"
	ctx1, cancelF1 := typex.NewCCTX()
	if err := engine.LoadOutEndWithCtx(SEMTECH_UDP_FORWARDER, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}
	defer cancelF1()
	grpcInend := typex.NewInEnd(typex.GRPC_SERVER,
		"rhilex Grpc InEnd",
		"rhilex Grpc InEnd", map[string]interface{}{
			"host": "127.0.0.1",
			"port": 2581,
		})
	grpcInend.UUID = "grpcInend"
	ctx, cancelF := typex.NewCCTX()
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Fatal(err)
	}

	rule := typex.NewRule(engine,
		"DATA_TO_SEMTECH_UDP_FORWARDER",
		"DATA_TO_SEMTECH_UDP_FORWARDER",
		"DATA_TO_SEMTECH_UDP_FORWARDER",
		grpcInend.UUID, "",
		`function Success() print("[Success]") end`,
		`
		Actions = {
			function(args)
				print("args --->",args)
				data:ToSemtechUdp('SEMTECH_UDP_FORWARDER', args)
				return true, args
			end
		}`,
		`function Failed(error) print("[Failed]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Fatal(err)
	}
	grpcConn, err := grpc.NewClient("127.0.0.1:2581",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer grpcConn.Close()
	client := rhilexrpc.NewRhilexRpcClient(grpcConn)
	for i := 0; i < 10; i++ {
		resp, err := client.Request(context.Background(), &rhilexrpc.RpcRequest{
			Value: fmt.Sprintf(`{"co2":10,"hum":30,"lex":22,"temp":100,"idx":%d}`, i),
		})
		if err != nil {
			t.Fatalf("grpc.Dial err: %v", err)
		}
		t.Logf("rhilex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(1 * time.Second)
	engine.Stop()
}
