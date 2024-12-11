package model

import (
	"database/sql/driver"
	"time"

	"encoding/json"
)

type RhilexModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"index"`
}
type StringList []string

/*
*
* 给GORM用的
*
 */
func (f StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *StringList) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

func (f StringList) String() string {
	b, _ := json.Marshal(f)
	return string(b)
}
func (f StringList) Len() int {
	return len(f)
}

type PageRequest struct {
	Current     int `json:"current,omitempty"`
	Size        int `json:"size,omitempty"`
	SearchCount int `json:"searchCount,omitempty"`
}

type PageResult struct {
	Current int `json:"current"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	Records any `json:"records"`
}
