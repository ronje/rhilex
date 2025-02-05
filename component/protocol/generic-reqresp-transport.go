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
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type Transport struct {
	buffer       []byte
	writer       io.Writer
	reader       io.Reader
	readTimeout  time.Duration
	writeTimeout time.Duration
	parser       GenericByteParser
	port         GenericPort
	logger       *logrus.Logger
}

func NewTransport(config ExchangeConfig) *Transport {
	return &Transport{
		port:         config.Port,
		logger:       config.Logger,
		readTimeout:  config.ReadTimeout,
		writeTimeout: config.WriteTimeout,
		reader:       bufio.NewReader(config.Port),
		writer:       bufio.NewWriter(config.Port),
		parser: GenericByteParser{
			edger:   config.PacketEdger,
			checker: &SimpleChecker{},
		},
		buffer: make([]byte, 0, 1024),
	}
}
func (transport *Transport) Write(data []byte) error {
	transport.port.SetWriteDeadline(time.Now().Add(
		transport.writeTimeout * time.Millisecond))
	defer transport.port.SetWriteDeadline(time.Time{})
	data = append(transport.parser.edger.Head[:], data...)
	data = append(data, transport.parser.edger.Tail[:]...)
	transport.logger.Debug("Transport.Write=", ByteDumpHexString(data))
	if _, err := transport.port.Write(data); err != nil {
		return err
	}

	return nil
}

// Read 方法从串口读取数据并解析
func (transport *Transport) Read() ([]byte, error) {
	transport.port.SetWriteDeadline(time.Now().Add(
		transport.readTimeout * time.Millisecond))
	defer transport.port.SetWriteDeadline(time.Time{})

	n, err := transport.port.Read(transport.buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read from serial port: %v", err)
	}

	if n == 0 {
		return nil, errors.New("no data read from serial port")
	}
	transport.logger.Debug("Transport.Read=", ByteDumpHexString(transport.buffer[:n]))
	packetData, parseErr := transport.parser.ParseBytes(transport.buffer[:n])
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse data: %v", parseErr)
	}

	return packetData, nil
}
func (transport *Transport) Status() error {
	if transport.port == nil {
		return fmt.Errorf("invalid port")
	}
	_, err := transport.port.Write([]byte{0})
	return err
}

func (transport *Transport) Close() error {
	return transport.port.Close()
}
