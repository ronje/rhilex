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
	"context"
	"io"
)

type GenericTransportSlaver struct {
	handler   *GenericProtocolHandler
	InChannel chan AppLayerFrame
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewGenericTransportSlaver(config TransporterConfig) *GenericTransportSlaver {
	ctx, cancel := context.WithCancel(context.Background())
	return &GenericTransportSlaver{
		ctx:       ctx,
		cancel:    cancel,
		handler:   NewGenericProtocolHandler(config),
		InChannel: make(chan AppLayerFrame, 1024),
	}
}

// Start
func (slaver *GenericTransportSlaver) StartLoop() {
	for {
		select {
		case <-slaver.ctx.Done():
			return
		default:
		}
		AppLayerFrame, errRead := slaver.handler.Read()
		if errRead != nil {
			if errRead == io.EOF {
				slaver.Stop()
			}
		}
		slaver.InChannel <- AppLayerFrame
	}
}

// Stop
func (slaver *GenericTransportSlaver) Stop() {
	if slaver.cancel != nil {
		slaver.cancel()
		slaver.handler.Close()
	}

}
