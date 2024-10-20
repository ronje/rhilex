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
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type TransportLayer struct {
	config TransporterConfig
}

func NewTransportLayer(config TransporterConfig) *TransportLayer {
	return &TransportLayer{config: config}
}

func (transport *TransportLayer) SendRequest(data []byte) error {
	if err := transport.config.Port.SetWriteDeadline(time.Now().Add(transport.config.WriteTimeout)); err != nil {
		return err
	}
	if _, err := transport.config.Port.Write(data); err != nil {
		return err
	}
	return transport.config.Port.SetWriteDeadline(time.Time{})
}

func (transport *TransportLayer) ReadResponse() ([]byte, error) {
	if err := transport.config.Port.SetReadDeadline(time.Now().Add(transport.config.ReadTimeout)); err != nil {
		return nil, err
	}
	responsetHeader := Header{}
	if err := binary.Read(transport.config.Port, binary.BigEndian, &responsetHeader); err != nil {
		return nil, err
	}
	responseLength := binary.BigEndian.Uint16(responsetHeader.Length[:])
	response := make([]byte, responseLength)
	if _, err := io.ReadFull(transport.config.Port, response); err != nil {
		return nil, err
	}

	if err := transport.config.Port.SetReadDeadline(time.Time{}); err != nil {
		return response, err
	}

	return response, nil
}

func (transport *TransportLayer) Status() error {
	if transport.config.Port == nil {
		return fmt.Errorf("invalid port")
	}
	_, err := transport.config.Port.Write([]byte{0})
	return err
}

func (transport *TransportLayer) Close() error {
	return transport.config.Port.Close()
}
