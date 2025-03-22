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

package protocol

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// 定义包头和包尾结构体
type PacketEdger struct {
	Head [2]byte
	Tail [2]byte
}
type DataChecker interface {
	CheckData(data []byte) error
}
type ExchangeConfig struct {
	Port         GenericPort
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PacketEdger  PacketEdger
	Logger       *logrus.Logger
}

func NewExchangeConfig() ExchangeConfig {
	return ExchangeConfig{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PacketEdger: PacketEdger{
			Head: [2]byte{0xAA, 0x55},
			Tail: [2]byte{0x0D, 0x0A},
		},
		Logger: logrus.New(),
	}
}

type GenericPort interface {
	io.ReadWriteCloser
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
