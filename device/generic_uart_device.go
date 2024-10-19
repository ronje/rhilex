package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/component/uartctrl"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type _UartCommonConfig struct {
	AutoRequest *bool `json:"autoRequest" validate:"required"`
}
type _UartRwConfig struct {
	Tag        string `json:"tag" validate:"required"`
	TimeSlice  uint64 `json:"timeSlice" validate:"required"`
	ReadFormat string `json:"readFormat" validate:"required" myself:"RAW,HEX,UTF8"` // 读取格式, "RAW"|"HEX"|"UTF8"
}
type _UartMainConfig struct {
	CommonConfig _UartCommonConfig `json:"commonConfig" validate:"required"`
	RwConfig     _UartRwConfig     `json:"rwConfig" validate:"required"`
	UartConfig   common.UartConfig `json:"uartConfig"`
}

type genericUartDevice struct {
	typex.XStatus
	serialPort serial.Port
	status     typex.DeviceState
	RuleEngine typex.Rhilex
	mainConfig _UartMainConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传，纯粹的串口读取网关
*
 */
func NewGenericUartDevice(e typex.Rhilex) typex.XDevice {
	uart := new(genericUartDevice)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = _UartMainConfig{
		CommonConfig: _UartCommonConfig{
			AutoRequest: func() *bool {
				b := true
				return &b
			}(),
		},
		RwConfig: _UartRwConfig{
			TimeSlice:  50,
			ReadFormat: "HEX",
			Tag:        "uart",
		},
		UartConfig: common.UartConfig{
			Timeout:  3000,
			Uart:     "/dev/ttyS1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *genericUartDevice) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId

	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	// CheckSerialBusy
	if err := uartctrl.CheckSerialBusy(uart.mainConfig.UartConfig.Uart); err != nil {
		return err
	}
	if uart.mainConfig.RwConfig.TimeSlice < 30 {
		errA := fmt.Errorf("TimeSlice Must Great than 30, but current is: %v",
			uart.mainConfig.RwConfig.TimeSlice)
		glogger.GLogger.Error(errA)
		return errA
	}
	ReadFormatTypes := []string{"HEX", "RAW", "UTF8"}
	if !slices.Contains(ReadFormatTypes, uart.mainConfig.RwConfig.ReadFormat) {
		errA := fmt.Errorf("ReadFormat Only Support Type: %v", ReadFormatTypes)
		glogger.GLogger.Error(errA)
		return errA
	}
	return nil
}

// 启动
func (uart *genericUartDevice) Start(cctx typex.CCTX) error {
	uart.Ctx = cctx.Ctx
	uart.CancelCTX = cctx.CancelCTX

	config := serial.Config{
		Address:  uart.mainConfig.UartConfig.Uart,
		BaudRate: uart.mainConfig.UartConfig.BaudRate,
		DataBits: uart.mainConfig.UartConfig.DataBits,
		Parity:   uart.mainConfig.UartConfig.Parity,
		StopBits: uart.mainConfig.UartConfig.StopBits,
		Timeout:  time.Duration(uart.mainConfig.UartConfig.Timeout) * time.Millisecond,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("serialPort start failed:", err)
		return err
	}

	uart.serialPort = serialPort
	if !*uart.mainConfig.CommonConfig.AutoRequest {
		uart.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		result := [2048]byte{}
		sliceTimer := time.NewTimer(time.Duration(uart.mainConfig.RwConfig.TimeSlice) * time.Millisecond)
		sliceTimer.Stop()
		peerCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-sliceTimer.C:
				mapV := map[string]interface{}{
					"tag": uart.mainConfig.RwConfig.Tag,
				}
				switch uart.mainConfig.RwConfig.ReadFormat {
				case "HEX":
					mapV["value"] = hex.EncodeToString(result[:peerCount])
				case "RAW":
					Value := []uint32{} // JSON会把[]Uint8识别为二进制，然后转换成Base64
					for i := 0; i < peerCount; i++ {
						Value = append(Value, uint32(result[i]))
					}
					mapV["value"] = Value
				case "UTF8":
					mapV["value"] = string(result[:peerCount])
				default:
					mapV["value"] = ""
					glogger.GLogger.Error("Not supported type:", uart.mainConfig.RwConfig.ReadFormat)
				}
				glogger.GLogger.Debug("Serial Port Read: ", result[:peerCount])
				bytes, _ := json.Marshal(mapV)
				uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
				for i := 0; i < peerCount; i++ {
					result[i] = 0 // 清空
				}
				peerCount = 0 // re-init index
			default:
				n, errR := io.ReadAtLeast(uart.serialPort, result[peerCount:], 1)
				if errR != nil {
					if !strings.Contains(errR.Error(), "timeout") {
						glogger.GLogger.Error(errR)
					}
				}
				if n != 0 {
					peerCount += n
					sliceTimer.Reset(time.Duration(uart.mainConfig.RwConfig.TimeSlice) * time.Millisecond)
				}
			}
		}
	}(uart.Ctx)
	uart.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来:
//
//	{
//	    "tag":"data tag",
//	    "value":"value s"
//	}
//
// t1.txt="OK"\xff\xff\xff
func (uart *genericUartDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	result := [2048]byte{}
	if string(cmd) == "HEX" {
		hexs, err1 := hex.DecodeString(string(cmd))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			return nil, err1
		}
		n, errSliceRequest := utils.SliceRequest(uart.Ctx, uart.serialPort,
			hexs, result[:], false, time.Duration(uart.mainConfig.RwConfig.TimeSlice)*time.Millisecond)
		if errSliceRequest != nil {
			return []byte{}, errSliceRequest
		}
		return result[:n], nil
	}
	if string(cmd) == "STRING" {
		// s := "t1.txt=\"RHILEX\"\xFF\xFF\xFF"
		n, err := uart.serialPort.Write(args)
		if err != nil {
			return nil, err
		}
		return result[:n], nil

	}
	return []byte{}, fmt.Errorf("unsupported cmd, must one of : STRING|HEX")
}

// 设备当前状态
func (uart *genericUartDevice) Status() typex.DeviceState {
	if uart.serialPort == nil {
		uart.status = typex.DEV_DOWN
	}
	return uart.status
}

// 停止设备
func (uart *genericUartDevice) Stop() {
	uart.status = typex.DEV_DOWN
	if uart.CancelCTX != nil {
		uart.CancelCTX()
	}
	if uart.serialPort != nil {
		uart.serialPort.Close()
	}

}

func (uart *genericUartDevice) Details() *typex.Device {
	return uart.RuleEngine.GetDevice(uart.PointId)
}

func (uart *genericUartDevice) SetState(status typex.DeviceState) {
	uart.status = status
}

func (uart *genericUartDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (uart *genericUartDevice) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (uart *genericUartDevice) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}
