package telemetry

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"gopkg.in/ini.v1"
)

func NewTelemetry() *Telemetry {
	return &Telemetry{
		MainConfig: TelemetryConfig{
			Enable:     false,
			ServerAddr: "127.0.0.1:9990",
		},
	}
}

type TelemetryConfig struct {
	Enable     bool   `ini:"enable"`
	ServerAddr string `ini:"server_addr"`
}

type Telemetry struct {
	MainConfig TelemetryConfig `json:"config"`
}

func (t *Telemetry) Init(section *ini.Section) error {

	if err := utils.InIMapToStruct(section, &t.MainConfig); err != nil {
		return err
	}
	if !t.MainConfig.Enable || len(t.MainConfig.ServerAddr) == 0 {
		return fmt.Errorf("Invalid config: %s", t.MainConfig.ServerAddr)
	}
	return nil
}
func (t *Telemetry) Start(typex.Rhilex) error {
	return sendMessage(t.MainConfig.ServerAddr)
}

/*
*
* 遥测的时候向服务器发送的数据
*
 */
type TelemetryInfo struct {
	Arch     string `json:"arch,omitempty"`
	OS       string `json:"os,omitempty"`
	StartAt  string `json:"start_at,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	Mac      string `json:"mac,omitempty"`
	Admin    string `json:"admin,omitempty"`
}

/*
*
* 发UDP包
*
 */
func sendMessage(ServerAddr string) error {
	conn, err := net.Dial("udp", ServerAddr)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	defer conn.Close()
	now := time.Now()
	formatted := now.Format("2006-01-02 15:04:05")
	info := TelemetryInfo{
		Admin:    typex.License.AuthorizeAdmin,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		StartAt:  formatted,
		Mac:      typex.License.MAC,
		DeviceId: typex.License.DeviceID,
	}
	data, _ := json.Marshal(&info)
	for i := 0; i < 5; i++ {
		_, err = conn.Write(data)
		if err != nil {
			glogger.GLogger.Error(err.Error())
			time.Sleep(300 * time.Millisecond)
			continue
		}
		break
	}
	return nil
}

func (t *Telemetry) Service(typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
func (t *Telemetry) Stop() error {
	return nil
}

func (t *Telemetry) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "BUSINESS_TELEMETRY",
		Name:        "Business Telemetry",
		Version:     "v0.0.1",
		Description: "Business Telemetry Statistics",
	}
}
