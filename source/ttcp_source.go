package source

import (
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TcpSourceConfig struct {
}
type TcpSource struct {
	typex.XStatus
	mainConfig TcpSourceConfig
	status     typex.SourceState
}

func NewTcpSource(e typex.Rhilex) typex.XSource {
	h := TcpSource{}
	h.RuleEngine = e
	return &h
}

func (hh *TcpSource) Init(inEndId string, configMap map[string]interface{}) error {
	hh.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &hh.mainConfig); err != nil {
		return err
	}
	return nil
}

func (hh *TcpSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	hh.status = typex.SOURCE_UP
	return nil
}

func (hh *TcpSource) Stop() {
	hh.status = typex.SOURCE_DOWN
	if hh.CancelCTX != nil {
		hh.CancelCTX()
	}
}

func (hh *TcpSource) Status() typex.SourceState {
	return hh.status
}

func (hh *TcpSource) Test(inEndId string) bool {
	return true
}

func (hh *TcpSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

// 来自外面的数据
func (*TcpSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*TcpSource) UpStream([]byte) (int, error) {
	return 0, nil
}
