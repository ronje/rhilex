package typex

//
//
// 创建资源的时候需要一个通用配置类
//
//

type XConfig struct {
	Type       string                `json:"type"` // 类型
	Engine     Rhilex                `json:"-"`
	NewDevice  func(Rhilex) XDevice  `json:"-"`
	NewSource  func(Rhilex) XSource  `json:"-"`
	NewTarget  func(Rhilex) XTarget  `json:"-"`
	NewCecolla func(Rhilex) XCecolla `json:"-"`
}
