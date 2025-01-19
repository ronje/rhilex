// Copyright (C) 2025 wwhai
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

package multimedia

import (
	"context"
	"fmt"
)

type BasicMultimediaStream struct {
	stream *MultimediaStream
	state  MultimediaStreamState
}

func (b *BasicMultimediaStream) Init(uuid string, configMap map[string]interface{}) error {
	// 模拟初始化逻辑
	fmt.Println("Initializing stream with UUID:", uuid)
	b.state = MEDIA_UP
	return nil
}

func (b *BasicMultimediaStream) Start(ctx context.Context) error {
	// 模拟启动逻辑
	fmt.Println("Starting stream")
	return nil
}

func (b *BasicMultimediaStream) OnCtrl(cmd []byte, args []byte) (any, error) {
	// 模拟控制逻辑
	fmt.Println("Received control command:", string(cmd))
	return nil, nil
}

func (b *BasicMultimediaStream) Status() MultimediaStreamState {
	return b.state
}

func (b *BasicMultimediaStream) Stop() {
	// 模拟停止逻辑
	b.state = MEDIA_STOP
	fmt.Println("Stopping stream")
}

func (b *BasicMultimediaStream) Details() *MultimediaStream {
	return b.stream
}
