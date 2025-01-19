// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package multimedia

import (
	"context"
)

type MultimediaStreamState int

const (
	// 故障
	MEDIA_DOWN MultimediaStreamState = 0
	// 启用
	MEDIA_UP MultimediaStreamState = 1
	// 暂停
	MEDIA_PAUSE MultimediaStreamState = 2
	// 停止
	MEDIA_STOP MultimediaStreamState = 3
	// 准备
	MEDIA_PENDING MultimediaStreamState = 4
	// 禁用
	MEDIA_DISABLE MultimediaStreamState = 5
)

type MultimediaStream struct {
	xMultimediaStream XMultimediaStream
	UUID              string
	Name              string
	Type              string
	Config            map[string]interface{}
	Description       string
}

// 多媒体流接口
type XMultimediaStream interface {
	Init(uuid string, configMap map[string]interface{}) error
	Start(context.Context) error
	OnCtrl(cmd []byte, args []byte) (any, error)
	Status() MultimediaStreamState
	Stop()
	Details() *MultimediaStream
}
