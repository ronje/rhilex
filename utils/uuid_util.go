package utils

import (
	"strings"

	"github.com/lithammer/shortuuid/v4"
)

// MakeUUID
func InUuid() string {
	return MakeUUID("IN")
}

// GoodsUuid
func GoodsUuid() string {
	return MakeUUID("GOODS")
}

// MakeUUID
func OutUuid() string {
	return MakeUUID("OUT")
}
func DeviceUuid() string {
	return MakeUUID("DEVICE")
}
func PluginUuid() string {
	return MakeUUID("PLUGIN")
}

func GroupUuid() string {
	return MakeUUID("GROUP")
}
func AppUuid() string {
	return MakeUUID("APP")
}
func AiBaseUuid() string {
	return MakeUUID("AIBASE")
}
func DataSchemaUuid() string {
	return MakeUUID("SCHEMA")
}

// MakeUUID
func RuleUuid() string {
	return MakeUUID("RULE")
}

// MakeUUID
func UserLuaUuid() string {
	return MakeUUID("USERLUA")
}

// MakeUUID
func MBusPointUUID() string {
	return MakeUUID("MBUS")
}

// MakeUUID
func ModbusPointUUID() string {
	return MakeUUID("MDTBUS")
}

// MakeUUID
func SiemensPointUUID() string {
	return MakeUUID("SIMENS")
}

// MakeUUID
func SnmpOidUUID() string {
	return MakeUUID("SNMPOID")
}

func BacnetPointUUID() string {
	return MakeUUID("BACNET")
}

func Dlt6452007PointUUID() string {
	return MakeUUID("DLT645")
}

func Cjt1882004PointUUID() string {
	return MakeUUID("CJT188")
}

func Szy2062016PointUUID() string {
	return MakeUUID("SZY206")
}

func UserProtocolPointUUID() string {
	return MakeUUID("USERPT")
}

func UnionPointUUID() string {
	return MakeUUID("UNION")
}

// MakeUUID
func MakeUUID(prefix string) string {
	return prefix + strings.ToUpper(shortuuid.New()[:8])
}
func MakeLongUUID(prefix string) string {
	return prefix + strings.ToUpper(shortuuid.New())
}
