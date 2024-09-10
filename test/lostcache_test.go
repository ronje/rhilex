package test

import (
	"context"

	"testing"
	"time"

	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// go test -timeout 30s -run ^Test_DataToGrepTime_Local_Cache github.com/hootrhino/rhilex/test -v -count=1

func Test_DataToGrepTime_Local_Cache(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "GRPC", "GRPC", map[string]interface{}{
		"port": 2581,
		"host": "127.0.0.1",
	})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		glogger.GLogger.Fatal("Rule load failed:", err)
	}

	OutEnd := typex.NewOutEnd(typex.GREPTIME_DATABASE,
		"GREPTIME_DATABASE", "GREPTIME_DATABASE", map[string]interface{}{
			"gwsn":     "rhilex",
			"host":     "127.0.0.1",
			"port":     4001,
			"username": "rhilex",
			"password": "rhilex",
			"database": "public",
			"table":    "public",
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
		glogger.GLogger.Fatal(err)
	}
	conn, err := grpc.NewClient("127.0.0.1:2581",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	defer conn.Close()

	client := rhilexrpc.NewRhilexRpcClient(conn)
	for i := 0; i < 10; i++ {
		resp, err := client.Work(context.Background(), &rhilexrpc.Data{
			Value: `{"co2":1.23,"hum":2.34,"lex":3.45,"temp":4.56}`,
		})
		if err != nil {
			glogger.GLogger.Fatal(err)
		}
		glogger.GLogger.Infof("ToGreptimeDB(%d)====: %v", i, resp.GetMessage())
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
	engine.Stop()
}
