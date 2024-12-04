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

package periphery

import (
	"golang.org/x/exp/slices"
)

/**
 * 获取厂商
 *
 */
func CheckVendor(env string) string {
	if slices.Contains([]string{
		"RHILEXG1",
		"RPI4B",
		"EN6400",
		"HAAS506LD1",
		"RHILEXPRO1",
	}, env) {
		return env
	}
	return "COMMON"
}
