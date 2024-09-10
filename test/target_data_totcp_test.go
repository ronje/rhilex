package test

import (
	"fmt"
	"log"
	"os/exec"

	"net"
	"testing"
	"time"

	httpserver "github.com/hootrhino/rhilex/component/apiserver"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/stretchr/testify/assert"

	"github.com/hootrhino/rhilex/typex"
)

var _DataToTcp_luaCase = `
function Main(arg)
	for i = 1, 3, 1 do
	local err = data:ToUdp('TcpServer',applib:T2J({temp = 20,humi = 13.45}))
	applib:log('result =>',err)
	time:Sleep(100)
	end
 return 0
end
`

// go test -timeout 30s -run ^Test_DataToTcp github.com/hootrhino/rhilex/test -v -count=1

func Test_DataToTcp(t *testing.T) {
	RmUnitTestDbFile(t)

	engine := RunTestEngine()
	engine.Start()

	hh := httpserver.NewHttpApiServer(engine)
	// UdpApiServer loaded default
	if err := engine.LoadPlugin("plugin.http_server", hh); err != nil {
		t.Fatal(err)
	}
	go _start_simple_Udp_server()
	//

	TcpServer := typex.NewOutEnd(typex.UDP_TARGET,
		"Udp", "Udp", map[string]interface{}{
			"host":             "127.0.0.1",
			"port":             8891,
			"cacheOfflineData": true,
		},
	)
	TcpServer.UUID = "TcpServer"
	ctx1, cancelF1 := typex.NewCCTX() // ,ctx, cancelF

	if err := engine.LoadOutEndWithCtx(TcpServer, ctx1, cancelF1); err != nil {
		t.Fatal(err)
	}

	uuid := _createTestTcpApp_1(t)
	time.Sleep(1 * time.Second)
	_updateTestTcpApp_1(t, uuid)

	time.Sleep(20 * time.Second)
	_deleteTestTcpApp_1(t, uuid)
	engine.Stop()
}

// --------------------------------------------------------------------------------------------------
// 起一个Udp服务器
// --------------------------------------------------------------------------------------------------
func _start_simple_TCP_server() {
	// TCP服务器监听
	tcp_addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8891")
	if err != nil {
		log.Fatal(err)
	}
	tcp_listener, err := net.ListenTCP("tcp", tcp_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer tcp_listener.Close()

	for {
		// 接受新的连接
		tcp_conn, err := tcp_listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		defer tcp_conn.Close()

		// 读取客户端数据
		data := make([]byte, 100)
		tcp_conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := tcp_conn.Read(data)
		if err != nil {
			log.Fatal(err)
		}
		tcp_conn.SetReadDeadline(time.Time{})

		// 打印接收的数据
		log.Println("TCP Received ============>:", string(data[:n]))
	}
}

//--------------------------------------------------------------------------------------------------
// 资源操作
//--------------------------------------------------------------------------------------------------

func _createTestTcpApp_1(t *testing.T) string {
	// 通过接口创建一个App
	body := `{"name": "testlua1","version": "1.0.0","autoStart": false,"description": "hello world"}`
	output, err := exec.Command("curl",
		"-X", "POST", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", body).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_createApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log(mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua1")
	assert.Equal(t, mApp[0].Version, "1.0.0")
	assert.Equal(t, mApp[0].AutoStart, false)
	return mApp[0].UUID
}
func _updateTestTcpApp_1(t *testing.T, uuid string) {
	body := `{"uuid": "%s","name": "testlua11","version": "1.0.1","autoStart": true,"luaSource":"AppNAME='OK1'\nAppVERSION='0.0.3'\n%s"}`

	t.Logf(body, uuid, _DataToTcp_luaCase)
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", fmt.Sprintf(body, uuid, _DataToTcp_luaCase)).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_updateApp: ", string(output))
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log("APP UUID ==> ", mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua11")
	assert.Equal(t, mApp[0].Version, "1.0.1")
	assert.Equal(t, mApp[0].AutoStart, true)
	// _startTestApp
	time.Sleep(1 * time.Second)
	_startTestTcpApp_1(t, mApp[0].UUID)
}
func _deleteTestTcpApp_1(t *testing.T, uuid string) {
	// 删除一个App
	output, err := exec.Command("curl", "-X", "DELETE", "http://127.0.0.1:2580/api/v1/app?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_deleteApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 0, len(mApp))
}
func _startTestTcpApp_1(t *testing.T, uuid string) {
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app/start?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_startTestApp: ", string(output))
}
