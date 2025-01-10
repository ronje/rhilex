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

package haas506

import "os"

func Init_HAAS506LD1() error {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		_HAAS506_AI_Init()
		_HAAS506_DI_Init()
		_HAAS506_DO_Init()
		_HAAS506_LED_Init()
		InitML307R4G(_ML307R_4G_PATH)
	}
	return nil
}
