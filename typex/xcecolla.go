// 抽象设备：
// 1.0 以后的大功能：支持抽象设备，抽象设备就是外挂的设备，Rhilex本来是个规则引擎，但是1.0之前的版本没有对硬件设备进行抽象支持
// 因此，1.0以后增加对硬件的抽象
// Target Source 描述了数据的流向，抽象设备描述了数据的载体。
// 举例：外挂一个设备，这个设备具备双工控制功能，例如电磁开关等，此时它强调的是设备的物理功能，而数据则不是主体。
// 因此需要抽象出来一个层专门来描述这些设备
package typex

type CecollaState int

const (
	// 故障
	CEC_DOWN CecollaState = 0
	// 启用
	CEC_UP CecollaState = 1
	// 暂停
	CEC_PAUSE CecollaState = 2
	// 停止
	CEC_STOP CecollaState = 3
	// 准备
	CEC_PENDING CecollaState = 4
	// 禁用
	CEC_DISABLE CecollaState = 5
)

type CecollaType string

func (d CecollaType) String() string {
	return string(d)

}

const (
	TENCENT_IOTHUB_CEC CecollaType = "TENCENT_IOTHUB_CEC" // 腾讯云物联网平台
	ITHINGS_IOTHUB_CEC CecollaType = "ITHINGS_IOTHUB_CEC" // ITHINGS物联网平台
)

type XCecolla interface {
	Init(CECId string, configMap map[string]interface{}) error
	Start(CCTX) error
	OnCtrl(cmd []byte, args []byte) ([]byte, error)
	Status() CecollaState
	Stop()
	Details() *Cecolla
	SetState(CecollaState)
}
