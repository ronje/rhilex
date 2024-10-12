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

package cjt1882004

import (
	"bytes"
	"fmt"
)

// CJT188DataReadRequest represents a CJT188 data read request.
type CJT188DataReadRequest struct {
	ControlCode byte   // Control code (e.g., 0x11 for read request)
	Length      byte   // Data length
	Address     []byte // Meter address
	DataID      byte   // Data identifier
}

// CJT188DataReadResponse represents a CJT188 data read response.
type CJT188DataReadResponse struct {
	ControlCode byte   // Control code (e.g., 0x91 for read response)
	Length      byte   // Data length
	Address     []byte // Meter address
	DataID      byte   // Data identifier
	Status      byte   // Status byte
	Data        []byte // Data field
	Checksum    byte   // Checksum
}

// PackCJT188DataReadRequest constructs a byte slice from a CJT188DataReadRequest.
func PackCJT188DataReadRequest(req CJT188DataReadRequest) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(req.ControlCode)
	buffer.WriteByte(req.Length)
	buffer.Write(req.Address)
	buffer.WriteByte(req.DataID)
	// Compute and append checksum
	checksum := calculateChecksum(buffer.Bytes())
	buffer.WriteByte(checksum)
	return buffer.Bytes()
}

// UnpackCJT188DataReadResponse parses a byte slice into a CJT188DataReadResponse.
func UnpackCJT188DataReadResponse(packet []byte) (*CJT188DataReadResponse, error) {
	if len(packet) < 8 { // Minimum frame length
		return nil, fmt.Errorf("invalid packet length")
	}

	var resp CJT188DataReadResponse
	resp.ControlCode = packet[0]
	resp.Length = packet[1]
	resp.Address = packet[2:6]
	resp.DataID = packet[6]
	resp.Status = packet[7]
	resp.Data = packet[8 : 8+resp.Length-3] // Subtract control code, length, and status

	// Checksum verification
	calculatedChecksum := calculateChecksum(packet[:len(packet)-1])
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
