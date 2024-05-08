package source

import (
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TemplateSourceConfig struct {
}
type TemplateSource struct {
	typex.XStatus
	mainConfig TemplateSourceConfig
	status     typex.SourceState
}

func NewTemplateSource(e typex.Rhilex) typex.XSource {
	h := TemplateSource{}
	h.RuleEngine = e
	return &h
}

func (hh *TemplateSource) Init(inEndId string, configMap map[string]interface{}) error {
	hh.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &hh.mainConfig); err != nil {
		return err
	}
	return nil
}

func (hh *TemplateSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	hh.status = typex.SOURCE_UP
	return nil
}

func (hh *TemplateSource) Stop() {
	hh.status = typex.SOURCE_DOWN
	if hh.CancelCTX != nil {
		hh.CancelCTX()
	}
}

func (hh *TemplateSource) Status() typex.SourceState {
	return hh.status
}

func (hh *TemplateSource) Test(inEndId string) bool {
	return true
}

func (hh *TemplateSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

// 来自外面的数据
func (*TemplateSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*TemplateSource) UpStream([]byte) (int, error) {
	return 0, nil
}
