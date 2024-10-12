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

import "fmt"

// MBusFrame represents the structure of an M-Bus frame.
type MBusFrame struct {
	ControlField byte
	Address      byte
	DataLength   byte
	Data         []byte
	Checksum     byte
}

// CalculateChecksum computes the checksum for the M-Bus frame.
func CalculateChecksum(frame MBusFrame) byte {
	var checksum byte
	checksum += frame.ControlField
	checksum += frame.Address
	checksum += frame.DataLength
	for _, b := range frame.Data {
		checksum += b
	}
	return ^checksum + 1 // Two's complement
}

// Pack constructs the M-Bus frame from the given data.
func (frame *MBusFrame) Pack() []byte {
	frame.Checksum = CalculateChecksum(*frame)
	packet := make([]byte, 0, 1+1+1+len(frame.Data)+1)
	packet = append(packet, frame.ControlField)
	packet = append(packet, frame.Address)
	packet = append(packet, frame.DataLength)
	packet = append(packet, frame.Data...)
	packet = append(packet, frame.Checksum)
	return packet
}

// Unpack parses the M-Bus frame from the given byte slice.
func Unpack(packet []byte) (*MBusFrame, error) {
	if len(packet) < 5 { // Minimum length of a valid M-Bus frame
		return nil, fmt.Errorf("invalid packet length")
	}

	frame := MBusFrame{
		ControlField: packet[0],
		Address:      packet[1],
		DataLength:   packet[2],
	}

	dataLength := int(frame.DataLength)
	if len(packet) < 5+dataLength { // Check if the packet contains enough data
		return nil, fmt.Errorf("packet too short for specified data length")
	}

	frame.Data = packet[3 : 3+dataLength]
	frame.Checksum = packet[3+dataLength]

	// Calculate and verify checksum
	calculatedChecksum := CalculateChecksum(frame)
	if frame.Checksum != calculatedChecksum {
		return nil, fmt.Errorf("invalid checksum")
	}

	return &frame, nil
}
