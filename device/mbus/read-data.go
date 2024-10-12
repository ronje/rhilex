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
	"bytes"
	"fmt"
)

// MBusDataReadRequest represents an M-Bus data read request.
type MBusDataReadRequest struct {
	Address byte // Slave address
	Command byte // Command (e.g., 0xFD for readout request)
}

// MBusDataReadResponse represents an M-Bus data read response.
type MBusDataReadResponse struct {
	Address    byte   // Slave address
	Command    byte   // Command (e.g., 0xFD for readout request)
	Status     byte   // Status byte
	DataLength byte   // Length of the data field
	Data       []byte // Data field
	Checksum   byte   // Checksum
}

// PackMBusDataReadRequest constructs a byte slice from an MBusDataReadRequest.
func PackMBusDataReadRequest(req MBusDataReadRequest) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x10) // Start of frame
	buffer.WriteByte(req.Address)
	buffer.WriteByte(req.Command)
	buffer.WriteByte(0x00) // Data length for request is always 0
	buffer.WriteByte(0x16) // End of frame
	return buffer.Bytes()
}

// UnpackMBusDataReadResponse parses a byte slice into an MBusDataReadResponse.
func UnpackMBusDataReadResponse(packet []byte) (*MBusDataReadResponse, error) {
	if len(packet) < 7 { // Minimum frame length
		return nil, fmt.Errorf("invalid packet length")
	}

	if packet[0] != 0x68 || packet[1] != 0x08 { // Check start of frame and control field
		return nil, fmt.Errorf("invalid frame header")
	}

	var resp MBusDataReadResponse
	resp.Address = packet[2]
	resp.Command = packet[3]
	resp.Status = packet[4]
	resp.DataLength = packet[5]
	resp.Data = packet[6 : 6+resp.DataLength]
	resp.Checksum = packet[6+resp.DataLength]

	// Checksum verification
	calculatedChecksum := calculateChecksum(packet[2 : 6+resp.DataLength])
	if resp.Checksum != calculatedChecksum {
		return nil, fmt.Errorf("invalid checksum")
	}

	return &resp, nil
}

// calculateChecksum computes the checksum for the given byte slice.
func calculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum += b
	}
	return checksum
}
