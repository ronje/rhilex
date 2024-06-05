// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package shellengine

import "syscall"

/*
*
* 系统调用
*
 */

const (
	iocNrBits   = 8
	iocTypeBits = 8
	iocSizeBits = 14
	iocDirBits  = 2

	iocNrShift   = 0
	iocTypeShift = iocNrShift + iocNrBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits

	iocNone  = 0
	iocWrite = 1
	iocRead  = 2
)

func IO(t, nr uintptr) uintptr {
	return IOC(iocNone, t, nr, 0)
}

func IOR(t, nr, size uintptr) uintptr {
	return IOC(iocRead, t, nr, size)
}

func IOW(t, nr, size uintptr) uintptr {
	return IOC(iocWrite, t, nr, size)
}

func IOWR(t, nr, size uintptr) uintptr {
	return IOC(iocRead|iocWrite, t, nr, size)
}

func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << iocDirShift) | (t << iocTypeShift) | (nr << iocNrShift) | (size << iocSizeShift)
}

func Ioctl(fd, request, value uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, request, value)
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

func IoctlX(fd, request, value uintptr) (int64, error) {
	x, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, request, value)
	if e == syscall.Errno(0) {
		return int64(x), nil
	}
	return 0, e
}
