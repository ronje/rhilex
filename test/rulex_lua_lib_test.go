package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hootrhino/rhilex/component/rhilexrpc"
	"github.com/hootrhino/rhilex/typex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_rhilex_base_lib(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	// Grpc Inend
	grpcInend := typex.NewInEnd("GRPC", "rhilex Grpc InEnd",
		"rhilex Grpc InEnd", map[string]any{
			"host": "127.0.0.1",
			"port": 2581,
		})
	ctx, cancelF := typex.NewCCTX() // ,ctx, cancelF
	if err := engine.LoadInEndWithCtx(grpcInend, ctx, cancelF); err != nil {
		t.Error("grpcInend load failed:", err)
	}
	//
	// Load Rule
	//
	rule := typex.NewRule(engine,
		"uuid4",
		"rule4",
		"rule4",
		[]string{grpcInend.UUID}[0],
		"",
		`function Success() print("[Success Callback]=> OK") end`,
		`
	Actions = {
		function(args)
			print("[rhilexlib:Time()] ==>", rhilexlib:Time())
			print("[rhilexlib:TsUnix()] ==>", rhilexlib:TsUnix())
			print("[rhilexlib:TsUnixNano()] ==>", rhilexlib:TsUnixNano())
			local MatchHexS = rhilexlib:MatchHex("age:[1,3];sex:[4,5]", "FFFFFF014CB2AA55")
			for key, value in pairs(MatchHexS) do
			    print('rhilexlib:MatchHex', key, value)
		    end
			-- rhilexlib:VSet('k', 'v')
			-- print("[rhilexlib:VGet(k)] ==>", rhilexlib:VGet('k'))
			-- Hello()
			-- rhilexlib:Throw('this is test Throw')
			return true, args
		end,
		function(args)
			rhilexlib:log(rhilexlib:Time())
			return true, args
		end
	}`,
		`function Failed(error) print("[Failed Callback]", error) end`)
	if err := engine.LoadRule(rule); err != nil {
		t.Error(err)
	}
	conn, err := grpc.Dial("127.0.0.1:2581", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
	}
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
	time.Sleep(5 * time.Second)
	engine.Stop()
}
