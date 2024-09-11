package test

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hootrhino/rhilex/typex"
)

// go test -timeout 30s -run ^Test_DataToUdp github.com/hootrhino/rhilex/test -v -count=1

func Test_DataToUdp(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// Grpc Inend
	grpcInend := typex.NewInEnd(typex.GRPC, "GRPC", "GRPC", map[string]interface{}{
		"port": 2581,
		"host": "127.0.0.1",
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Error("Rule load failed:", err)
	}

	OutEnd := typex.NewOutEnd(typex.UDP_TARGET,
		"UDP_TARGET", "UDP_TARGET", map[string]interface{}{
			"host":             "127.0.0.1",
			"port":             8891,
			"pingPacket":       "rhilex\r\n",
			"cacheOfflineData": true,
			"allowPing":        true,
			"dataMode":         "RAW_STRING",
			"timeout":          3000,
		},
	)
	OutEnd.UUID = "Test"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(OutEnd, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}
	rule := typex.NewRule(engine, "Test", "Test", "Test",
		OutEnd.UUID,
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
		t.Error(err)
	}
	conn, err := grpc.NewClient("127.0.0.1:2581",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
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
