package test

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// go test -timeout 30s -run ^Test_DataToGrepTime github.com/hootrhino/rhilex/test -v -count=1

func Test_DataToGrepTime(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// Grpc Inend
	grpcInend := typex.NewInEnd(typex.GRPC, "GRPC", "GRPC", map[string]interface{}{
		"port": 2581,
		"host": "127.0.0.1",
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Error("Rule load failed:", err)
	}

	OutEnd := typex.NewOutEnd(typex.GREPTIME_DATABASE,
		"GREPTIME_DATABASE", "GREPTIME_DATABASE", map[string]interface{}{
			"gwsn":             "rhilex",
			"host":             "127.0.0.1",
			"port":             4001,
			"username":         "rhilex",
			"password":         "rhilex",
			"database":         "public",
			"table":            "public",
			"cacheOfflineData": true,
		})
	OutEnd.UUID = "Test"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(OutEnd, ctx1, cancelF1); err != nil {
		glogger.GLogger.Fatal(err)
	}
	rule := typex.NewRule(engine, "Test", "Test", "Test",
		grpcInend.UUID,
		"",
		`function Success() end`,
		`
		Actions = {
			function(args)
			    local msg, err = data:ToGreptimeDB('Test', args)
				print("[data To GrepTimeDb] ============= ", msg, err)
				return true, args
			end
		}`,
		`function Failed(error) end`)
	if err := engine.LoadRule(rule); err != nil {
		glogger.GLogger.Error(err)
	}
	conn, err := grpc.NewClient("127.0.0.1:2581",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			t.Fatalf("grpc.Dial err: %v", err)
		}
		t.Logf("rhilex Rpc Call Result ====>>: %v", resp.GetMessage())
	}

	time.Sleep(2 * time.Second)
	engine.Stop()
}
