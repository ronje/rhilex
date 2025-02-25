//go:build windows
// +build windows

// Copyright (C) 2025 wwhai
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

package en6400

// AccelerationData 存储加速度数据
type AccelerationData struct {
	X float64
	Y float64
	Z float64
}

// ReadAcceleration 读取加速度数据，在 Windows 下返回占位数据
func ReadAcceleration() AccelerationData {
	return AccelerationData{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}
}
