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

package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
)
import "testing"

// go test -timeout 30s -run ^TestCheckFileType github.com/hootrhino/rhilex/test -v -count=1
func TestCheckFileType(t *testing.T) {
	t.Log(CheckFileType("../rhilex"))
	t.Log(CheckFileType("../rhilex-arm32linux"))
	t.Log(CheckFileType("../rhilex.exe"))
}
func CheckFileType(filePath string) error {
	currentArch := runtime.GOARCH
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var magicNumber [4]byte
	_, err = file.Read(magicNumber[:])
	if err != nil {
		return err
	}
	switch {
	// ELF文件，用于Linux
	case bytes.Equal(magicNumber[:], []byte{0x7F, 'E', 'L', 'F'}):
		elfArch, err := checkELFArch(file)
		if err != nil {
			return err
		}
		if elfArch != currentArch {
			return fmt.Errorf("ELF architecture mismatch: %s != %s", elfArch, currentArch)
		}
		// Windows
	case CheckPEFileMagic(magicNumber):
		return fmt.Errorf("not support windows PE")
	case CheckDOSHeaderMagic(magicNumber):
		return fmt.Errorf("not support windows DOS")
	default:
		return fmt.Errorf("unknown file type")
	}

	return nil
}

// checkELFArch 检查ELF文件的架构
func checkELFArch(file *os.File) (string, error) {
	type elfHeader struct {
		Ident     [16]byte
		Type      uint16
		Machine   uint16
		Version   uint32
		Entry     uint64
		Phoff     uint64
		Shoff     uint64
		Flags     uint32
		Ehsize    uint16
		Phentsize uint16
		Phnum     uint16
		Shentsize uint16
		Shnum     uint16
		Shstrndx  uint16
	}
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	var hdr elfHeader
	err = binary.Read(file, binary.LittleEndian, &hdr)
	if err != nil {
		return "", err
	}
	switch hdr.Machine {
	case 3:
		return "386", nil // x86
	case 62:
		return "amd64", nil // x86_64
	case 40:
		return "arm", nil // ARM

	default:
		return "", fmt.Errorf("unknown ELF architecture")
	}
}

func CheckPEFileMagic(data [4]byte) bool {
	return (uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24) == 0x50450000
}

func CheckDOSHeaderMagic(data [4]byte) bool {
	return (uint32(data[0]) | uint32(data[1])<<8) == 0x5A4D
}
