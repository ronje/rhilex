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

package mbus

import (
	"bufio"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/sirupsen/logrus"
)

type MbusClientHandler struct {
	logger        *logrus.Logger
	DataLinkLayer *MbusDataLinkLayer
	Transporter   *MbusSerialTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewMbusClientHandler(Transporter typex.GenericRWC) *MbusClientHandler {
	handler := &MbusClientHandler{
		DataLinkLayer: &MbusDataLinkLayer{},
		Transporter: NewMbusSerialTransporter(TransporterConfig{
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
			Port:         Transporter,
		}),
	}
	return handler
}

/**
 * 关闭
 *
 */
func (handler *MbusClientHandler) Close() error {
	if handler.Transporter.port == nil {
		return fmt.Errorf("invalid Transporter")
	}
	return handler.Transporter.port.Close()
}

func (handler *MbusClientHandler) SetLogger(logger *logrus.Logger) {
	handler.logger = logger
}

/**
 * 发送请求
 *
 */
func (handler *MbusClientHandler) Request(data []byte) ([]byte, error) {
	request := ""
	for _, b := range data {
		request += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("MbusClientHandler.Request:", request)
	r, e := handler.Transporter.SendFrame(data)
	handler.logger.Debug("MbusClientHandler Transporter.Received:", utils.ByteDumpHexString(r))
	return r, e
}

func (handler *MbusClientHandler) RequestUD2(primaryID uint8) ShortFrame {
	data := NewShortFrame()
	data[1] = 0x5b
	data[2] = primaryID
	data.SetChecksum()
	return data
}

// SndNKE slave will ack with SingleCharacterFrame (e5).
func (handler *MbusClientHandler) SndNKE(primaryID uint8) ShortFrame {
	data := NewShortFrame()
	data[1] = 0x40
	data[2] = primaryID
	data.SetChecksum()
	return data
}

func (handler *MbusClientHandler) ApplicationReset(primaryID uint8) LongFrame {
	data := LongFrame{
		0x68, // Start byte long/control
		0x06, // length
		0x06, // length
		0x68, // Start byte long/control
		0x73, // SND_UD
		primaryID,
		0x50, // CI field data send
		0x00, // checksum
		0x16, // stop byte
	}

	data.SetLength()
	data.SetChecksum()
	return data
}

// water meter has 19004636 7704 14 07.
func (handler *MbusClientHandler) SendUD2() LongFrame {
	data := LongFrame{
		0x68, // Start byte long/control
		0x00, // length
		0x00, // length
		0x68, // Start byte long/control

		0x73, // REQ_UD2
		0xFD,
		0x52, // CI-field selection of slave

		0x00, // address
		0x00, // address
		0x00, // address
		0x00, // address

		0xFF, // manufacturer code
		0xFF, // manufacturer code

		0xFF, // id

		0xFF, // medium code

		0x00, // checksum
		0x16, // stop byte
	}

	data.SetLength()
	data.SetChecksum()

	return data
}

func (handler *MbusClientHandler) SetPrimaryUsingSecondary(secondary uint64, primary uint8) LongFrame {
	data := LongFrame{
		0x68, // Start byte long/control
		0x00, // length
		0x00, // length
		0x68, // Start byte long/control
		0x73, // SND_UD
		0xFD,
		0x51, // CI field data send
		0x00, // address
		0x00, // address
		0x00, // address
		0x00, // address
		0xFF, // manufacturer code
		0xFF, // manufacturer code
		0xFF, // id
		0xFF, // medium code
		0x01, // DIF field
		0x7a, // VIF field
		primary,
		0x00, // checksum
		0x16, // stop byte
	}

	a := uintToBCD(secondary, 4)
	data[7] = a[0]
	data[8] = a[1]
	data[9] = a[2]
	data[10] = a[3]

	data.SetLength()
	data.SetChecksum()
	return data
}

func (handler *MbusClientHandler) SetPrimaryUsingPrimary(oldPrimary uint8, newPrimary uint8) LongFrame {
	data := LongFrame{
		0x68, // Start byte long/control
		0x06, // length
		0x06, // length
		0x68, // Start byte long/control
		0x73, // REQ_UD2
		oldPrimary,
		0x51, // CI field data send
		0x01, // DIF field
		0x7a, // VIF field
		newPrimary,
		0x00, // checksum
		0x16, // stop byte
	}

	data.SetLength()
	data.SetChecksum()
	return data
}

// ReadAllFrames supports FCB and reads out all frames from the device using primaryID.
func (handler *MbusClientHandler) ReadAllFrames(primaryID int) ([]*DecodedFrame, error) {
	frame := handler.SndNKE(uint8(primaryID))
	handler.logger.Debugf("sending nke: % x\n", frame)
	_, err := handler.Transporter.port.Write(frame)
	if err != nil {
		return nil, err
	}

	_, err = handler.ReadSingleCharFrame()
	if err != nil {
		return nil, err
	}

	frames := []*DecodedFrame{}
	respFrame := &DecodedFrame{}
	lastFCB := true
	frameCnt := 0
	for respFrame.HasMoreRecords() || frameCnt == 0 {
		frameCnt++
		frame = handler.RequestUD2(uint8(primaryID))
		if !lastFCB {
			frame.SetFCB()
			frame.SetChecksum()
		}
		lastFCB = frame.C().FCB()

		_, err = handler.Transporter.port.Write(frame)
		if err != nil {
			return nil, err
		}

		resp, err := handler.ReadLongFrame()
		if err != nil {
			return nil, err
		}

		respFrame, err = resp.Decode()
		if err != nil {
			return nil, err
		}
		frames = append(frames, respFrame)
	}

	return frames, nil
}

// ReadSingleFrame reads one frame from the device. Does not reset device before asking.
func (handler *MbusClientHandler) ReadSingleFrame(primaryID int) (*DecodedFrame, error) {
	frame := handler.RequestUD2(uint8(primaryID))
	if _, err := handler.Transporter.port.Write(frame); err != nil {
		return nil, err
	}

	resp, err := handler.ReadLongFrame()
	if err != nil {
		return nil, err
	}

	respFrame, err := resp.Decode()
	if err != nil {
		return nil, err
	}

	return respFrame, nil
}

var ErrNoLongFrameFound = fmt.Errorf("no long frame found")

func (handler *MbusClientHandler) ReadLongFrame() (LongFrame, error) {
	buf := make([]byte, 4096)
	tmp := make([]byte, 4096)

	// foundStart := false
	length := 0
	globalN := -1
	for {
		n, err := handler.Transporter.port.Read(tmp)
		if err != nil {
			return LongFrame{}, fmt.Errorf("error reading from connection: %w", err)
		}

		for _, b := range tmp[:n] {
			globalN++
			buf[globalN] = b

			if globalN > 256 {
				return LongFrame{}, ErrNoLongFrameFound
			}

			// look for end byte after length +C+A+CI+checksum
			if length != 0 && globalN > length+4 && b == 0x16 {
				return LongFrame(buf[:globalN+1]), nil
			}

			// look for start sequence 68 LL LL 68
			if length == 0 && buf[0] == 0x68 && buf[3] == 0x68 && buf[1] == buf[2] {
				length = int(buf[1])
			}
		}
	}
}

func (handler *MbusClientHandler) ReadSingleCharFrame() (LongFrame, error) {
	buf := bufio.NewReader(handler.Transporter.port)
	msg, err := buf.ReadBytes(SingleCharacterFrame)
	if err != nil {
		return nil, err
	}
	return LongFrame(msg), nil
}

// ReadAnyAndPrint is used for debugging.
func (handler *MbusClientHandler) ReadAnyAndPrint() error {
	tmp := make([]byte, 256) // using small tmo buffer for demonstrating
	for {
		n, err := handler.Transporter.port.Read(tmp)
		if err != nil {
			return err
		}
		handler.logger.Debugf("% x\n", tmp[:n])
	}
}
