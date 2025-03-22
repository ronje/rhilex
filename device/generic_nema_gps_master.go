package device

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type NemaGpsConfig struct {
	Parse *bool  `json:"parse" validate:"required"`
	GwSN  string `json:"gwsn" validate:"required"`
}

type NemaGpsMainConfig struct {
	GpsConfig  NemaGpsConfig        `json:"gpsConfig" validate:"required"`
	UartConfig resconfig.UartConfig `json:"uartConfig" validate:"required"`
}

type NemaGpsMasterDevice struct {
	typex.XStatus
	serialPort serial.Port
	status     typex.SourceState
	RuleEngine typex.Rhilex
	mainConfig NemaGpsMainConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传，纯粹的串口读取网关
*
 */
func NewNemaGpsMasterDevice(e typex.Rhilex) typex.XDevice {
	gpsd := new(NemaGpsMasterDevice)
	gpsd.locker = &sync.Mutex{}
	gpsd.mainConfig = NemaGpsMainConfig{
		GpsConfig: NemaGpsConfig{
			Parse: new(bool),
			GwSN:  "rhilex",
		},
		UartConfig: resconfig.UartConfig{
			Timeout:  3000,
			Uart:     "/dev/ttyS1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	gpsd.RuleEngine = e
	return gpsd
}

//  初始化
func (gpsd *NemaGpsMasterDevice) Init(devId string, configMap map[string]any) error {
	gpsd.PointId = devId
	intercache.RegisterSlot(gpsd.PointId)

	if err := utils.BindSourceConfig(configMap, &gpsd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// 启动
func (gpsd *NemaGpsMasterDevice) Start(cctx typex.CCTX) error {
	gpsd.Ctx = cctx.Ctx
	gpsd.CancelCTX = cctx.CancelCTX

	config := serial.Config{
		Address:  gpsd.mainConfig.UartConfig.Uart,
		BaudRate: gpsd.mainConfig.UartConfig.BaudRate,
		DataBits: gpsd.mainConfig.UartConfig.DataBits,
		Parity:   gpsd.mainConfig.UartConfig.Parity,
		StopBits: gpsd.mainConfig.UartConfig.StopBits,
		Timeout:  time.Duration(gpsd.mainConfig.UartConfig.Timeout) * time.Millisecond,
	}
	serialPort, errOpen := serial.Open(&config)
	if errOpen != nil {
		glogger.GLogger.Error("serial port start failed err:", errOpen, ", config:", config)
		return errOpen
	}

	gpsd.serialPort = serialPort
	go func(ctx context.Context) {
		defer gpsd.serialPort.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			rawAiSString, err := readUntilNewLine(gpsd.serialPort)
			if err != nil {
				glogger.Error(err)
				continue
			}
			if !*gpsd.mainConfig.GpsConfig.Parse {
				ds := `{"gwsn":"%s","data":"%s"}`
				lens := len(rawAiSString)
				if lens > 2 {
					gpsd.RuleEngine.WorkDevice(gpsd.Details(),
						fmt.Sprintf(ds, gpsd.mainConfig.GpsConfig.GwSN, rawAiSString))
				}
			} else {
				GPGGAData, errParse := parseGPGGA(rawAiSString)
				if errParse != nil {
					glogger.GLogger.Error("parse GPGGA error:", errParse)
					continue
				}
				gpsd.RuleEngine.WorkDevice(gpsd.Details(), GPGGAData.String())
			}
		}
	}(gpsd.Ctx)
	gpsd.status = typex.SOURCE_UP
	return nil
}

func (gpsd *NemaGpsMasterDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, fmt.Errorf("unsupported cmd")
}

// 设备当前状态
func (gpsd *NemaGpsMasterDevice) Status() typex.SourceState {
	if gpsd.serialPort == nil {
		gpsd.status = typex.SOURCE_DOWN
	}
	return gpsd.status
}

// 停止设备
func (gpsd *NemaGpsMasterDevice) Stop() {
	gpsd.status = typex.SOURCE_DOWN
	intercache.UnRegisterSlot(gpsd.PointId)
	if gpsd.CancelCTX != nil {
		gpsd.CancelCTX()
	}
	if gpsd.serialPort != nil {
		gpsd.serialPort.Close()
	}

}

func (gpsd *NemaGpsMasterDevice) Details() *typex.Device {
	return gpsd.RuleEngine.GetDevice(gpsd.PointId)
}

func (gpsd *NemaGpsMasterDevice) SetState(status typex.SourceState) {
	gpsd.status = status
}

func (gpsd *NemaGpsMasterDevice) OnDCACall(UUID string, Command string, Args any) typex.DCAResult {
	return typex.DCAResult{}
}

// GPGGAData holds the parsed data from a GPGGA NMEA sentence
type GPGGAData struct {
	UTC        string
	Latitude   float64
	Longitude  float64
	Quality    int
	Satellites int
	HDOP       float64
	Altitude   float64
}

func (O GPGGAData) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

// parseGPGGA parses a GPGGA NMEA sentence and returns a GPGGAData struct
func parseGPGGA(sentence string) (*GPGGAData, error) {
	// Check if the sentence starts with $GPGGA
	if !strings.HasPrefix(sentence, "$GPGGA") {
		return nil, fmt.Errorf("invalid GPGGA sentence")
	}

	// Split the sentence into fields
	fields := strings.Split(sentence, ",")

	// Check if the sentence has the correct number of fields
	if len(fields) < 14 {
		return nil, fmt.Errorf("invalid number of fields in GPGGA sentence")
	}

	// Parse the UTC time
	utc := fields[1]

	// Parse the latitude
	lat, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %v", err)
	}
	latMinutes := math.Mod(lat, 100)
	lat = (lat / 100) + (latMinutes / 60)

	// Adjust the latitude for hemisphere
	if fields[3] == "S" {
		lat = -lat
	}

	// Parse the longitude
	lon, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %v", err)
	}
	lonMinutes := math.Mod(lon, 100)
	lon = (lon / 100) + (lonMinutes / 60)

	// Adjust the longitude for hemisphere
	if fields[5] == "W" {
		lon = -lon
	}

	// Parse the quality
	quality, err := strconv.Atoi(fields[6])
	if err != nil {
		return nil, fmt.Errorf("invalid quality: %v", err)
	}

	// Parse the number of satellites
	satellites, err := strconv.Atoi(fields[7])
	if err != nil {
		return nil, fmt.Errorf("invalid number of satellites: %v", err)
	}

	// Parse the HDOP
	hdop, err := strconv.ParseFloat(fields[8], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid HDOP: %v", err)
	}

	// Parse the altitude
	altitude, err := strconv.ParseFloat(fields[9], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid altitude: %v", err)
	}

	// Create and populate the GPGGAData struct
	gpggaData := &GPGGAData{
		UTC:        utc,
		Latitude:   lat,
		Longitude:  lon,
		Quality:    quality,
		Satellites: satellites,
		HDOP:       hdop,
		Altitude:   altitude,
	}

	return gpggaData, nil
}

// readUntilNewLine reads from the provided reader until a newline character is encountered.
// It handles both \n and \r\n as newline characters.
// It returns the line read (excluding the newline characters) and any error encountered.
func readUntilNewLine(reader io.Reader) (string, error) {
	var buffer bytes.Buffer
	for {
		var b [1]byte
		_, err := reader.Read(b[:])
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					// Return the remaining content if we hit EOF
					return buffer.String(), nil
				}
				break // End of file reached with no content
			}
			return "", err
		}
		if b[0] == '\n' {
			break
		}
		if b[0] == '\r' {
			// Check if the next character is a newline
			var nextByte [1]byte
			_, err := reader.Read(nextByte[:])
			if err != nil {
				if err == io.EOF {
					// EOF after a carriage return is treated as a newline
					break
				}
				return "", err
			}
			if nextByte[0] == '\n' {
				// Skip the newline character
				continue
			}
			// If the next character is not a newline, put the carriage return back
			buffer.WriteByte(b[0])
		}
		buffer.WriteByte(b[0])
	}
	return buffer.String(), nil
}
