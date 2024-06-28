// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package transceivercom

import (
	"encoding/json"
	"time"
)

type TransceiverConfig map[string]any

func (O TransceiverConfig) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

type TransceiverType uint8

const (
	COMMON_RF TransceiverType = 0
	WIFI      TransceiverType = 1
	BLC       TransceiverType = 2
	BLE       TransceiverType = 3
	ZIGBEE    TransceiverType = 4
	RF24g     TransceiverType = 5
	RF433M    TransceiverType = 6
	MN4G      TransceiverType = 7
	MN5G      TransceiverType = 8
	NBIoT     TransceiverType = 9
	LORA      TransceiverType = 10
	LORA_WAN  TransceiverType = 11
)

type TransceiverStatusCode uint8

const (
	TC_ERROR TransceiverStatusCode = 0
	TC_UP    TransceiverStatusCode = 1
	TC_DOWN  TransceiverStatusCode = 2
)

type CommunicatorInfo struct {
	Name     string          `json:"name"`
	Model    string          `json:"model"`
	Mac      string          `json:"mac"`
	Firmware string          `json:"firmware"`
	Type     TransceiverType `json:"type"`
	Vendor   string          `json:"vendor"`
}

func (O CommunicatorInfo) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

type TransceiverStatus struct {
	Code  TransceiverStatusCode
	Error error
}

type TransceiverCommunicator interface {
	Start(TransceiverConfig) error
	Ctrl(cmd []byte, timeout time.Duration) ([]byte, error)
	Status() TransceiverStatus
	Info() CommunicatorInfo
	Stop()
}
