package dto

type BacnetRouterDataPointVO struct {
	UUID           string `json:"uuid"`
	DeviceUUID     string `json:"device_uuid"`
	Tag            string `json:"tag"`
	Alias          string `json:"alias"`
	BacnetDeviceId uint32 `json:"bacnetDeviceId"`
	ObjectType     string `json:"objectType"`
	ObjectId       uint32 `json:"objectId"`
	ErrMsg         string `json:"errMsg"`        // 运行时数据
	Status         uint32 `json:"status"`        // 运行时数据
	LastFetchTime  uint64 `json:"lastFetchTime"` // 运行时数据
	Value          string `json:"value"`         // 运行时数据
}

type BacnetRouterDataPointCreateOrUpdate struct {
	UUID       string `json:"uuid"`
	Tag        string `json:"tag"`
	Alias      string `json:"alias"`
	ObjectType string `json:"objectType"`
	ObjectId   uint32 `json:"objectId"`
}
