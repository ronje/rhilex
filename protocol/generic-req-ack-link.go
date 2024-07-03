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
	"fmt"
	"io"
)

const MaxBufferSize = 1024

// Send sends the data to the specified channel.
func Send(channel chan<- []byte, data []byte) {
	channel <- data
}

// ReadFromIO reads from the io.ReadWriteCloser and sends complete messages to the channel.
// Messages are considered complete when they match the pattern: EE EF [DATA] \r \n
// ______--______--______--
func ReadFromIO(ctx context.Context, ioRw io.ReadWriteCloser, channel chan<- []byte) error {
	buffer := make([]byte, MaxBufferSize)
	var acc int
	edgeSignal1 := false
	edgeSignal2 := false

	current := [1]byte{}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		default:
		}
		_, errRead := ioRw.Read(current[:])
		if errRead != nil {
			if errRead == io.EOF {
				continue
			}
			return errRead
		}
		if acc > 0 {
			if current[0] == '\xEF' && buffer[acc-1] == '\xEE' {
				edgeSignal1 = true
			}
			if current[0] == '\n' && buffer[acc-1] == '\r' {
				edgeSignal2 = true
			}
		}

		if edgeSignal1 && edgeSignal2 {
			Send(channel, buffer[:acc])
			acc = 0
			edgeSignal1 = false
			edgeSignal2 = false
		} else {
			buffer[acc] = current[0]
			acc++
			if acc >= MaxBufferSize {
				acc = 0
				edgeSignal1 = false
				edgeSignal2 = false
			}
		}
	}
}

// ReadFromChannel reads data from an input channel and processes it according to specific rules,
// then sends the processed data to an output channel until the context is canceled.
// It checks for specific sequences in the input data and sends the accumulated buffer to the output
// channel when a sequence is detected.
func ReadFromChannel(ctx context.Context, inChannel chan []byte) ([]byte, error) {
	buffer := make([]byte, MaxBufferSize)
	var acc int
	edgeSignal1 := false
	edgeSignal2 := false
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled")
		case BinData := <-inChannel:
			for _, currentByte := range BinData {
				buffer[acc] = currentByte
				if acc > 0 {
					if currentByte == '\xEF' && buffer[acc-1] == '\xEE' {
						edgeSignal1 = true
					}
					if currentByte == '\n' && buffer[acc-1] == '\r' {
						edgeSignal2 = true
					}
				}
				if edgeSignal1 && edgeSignal2 {
					acc = 0
					edgeSignal1 = false
					edgeSignal2 = false
					return buffer[:acc], nil
				}
				acc++
				if acc >= MaxBufferSize {
					acc = 0
					edgeSignal1 = false
					edgeSignal2 = false
				}
			}
		}
	}
}
