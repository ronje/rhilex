package dto

type DataPointVO struct {
	UUID          string                 `json:"uuid"`
	DeviceUUID    string                 `json:"device_uuid"`
	Tag           string                 `json:"tag"`
	Alias         string                 `json:"alias"`
	Frequency     int                    `json:"frequency"`
	Config        map[string]interface{} `json:"config"`
	ErrMsg        string                 `json:"errMsg"`        // 运行时数据
	Status        uint32                 `json:"status"`        // 运行时数据
	LastFetchTime uint64                 `json:"lastFetchTime"` // 运行时数据
	Value         string                 `json:"value"`         // 运行时数据
}
