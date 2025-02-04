package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rhilex/component/apiserver"
	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/plugin"
	"github.com/hootrhino/rhilex/typex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* Test_data_to_tdengine
*
 */
func Test_data_to_tdengine(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)

	// HttpApiServer loaded default
	if err := plugin.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	// Grpc Inend
	grpcInend := typex.NewInEnd(
		"GRPC",
		"Test_data_to_tdengine",
		"Test_data_to_tdengine", map[string]interface{}{
			"port":             2581,
			"host":             "127.0.0.1",
			"cacheOfflineData": true,
		})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Fatal("grpcInend load failed:", err)
	}

	tdOutEnd := typex.NewOutEnd(typex.TDENGINE_TARGET,
		"Test_data_to_tdengine",
		"Test_data_to_tdengine",
		map[string]interface{}{
			"fqdn":           "127.0.0.1",
			"port":           6041,
			"username":       "root",
			"password":       "taosdata",
			"dbName":         "device",
			"createDbSql":    "CREATE DATABASE IF NOT EXISTS device UPDATE 0;",
			"createTableSql": "CREATE TABLE IF NOT EXISTS meter01 (ts TIMESTAMP, co2 INT, hum INT, lex INT, temp INT);",
			"insertSql":      "INSERT INTO meter01 VALUES (NOW, %v, %v, %v, %v);",
		})
	tdOutEnd.UUID = "TD1"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadOutEndWithCtx(tdOutEnd, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}
	//
	// Load Rule [{"co2":10,"hum":30,"lex":22,"temp":100}]
	//
	callback :=
		`Actions = {
			function(args)
				local t = rhilexlib:J2T(data)
				local Result = data:ToTdEngine('TD1', string.format("%d, %d, %d, %d", t['co2'], t['hum'], t['lex'], t['temp']))
				print("data:ToTdEngine Result", Result==nil)
				return false, data
			end
		}`
	rule1 := typex.NewRule(engine,
		"uuid1",
		"rule1",
		"rule1",
		[]string{grpcInend.UUID}[0],
		"",
		`function Success() print("[Test_data_to_tdengine Success Callback]=> OK") end`,
		callback,
		`function Failed(error) print("[Test_data_to_tdengine Failed Callback]", error) end`)

	if err := engine.LoadRule(rule1); err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := rhilexrpc.NewRhilexRpcClient(conn)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	rng.Int()
	for i := 0; i < 10; i++ {
		resp, err := client.Request(context.Background(), &rhilexrpc.RpcRequest{
			Value: fmt.Sprintf(`{"co2":10,"hum":30,"lex":22,"temp":100,"idx":%d}`, i),
		})
		if err != nil {
			t.Fatalf("grpc.Dial err: %v", err)
		}
		t.Logf("rhilex Rpc Call Result ====>>: %v", resp.GetMessage())
	}
	time.Sleep(3 * time.Second)
	engine.Stop()
}
