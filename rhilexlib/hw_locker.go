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

package rhilexlib

// 资源锁，防止同时操作硬件的时候引发竞争问题，比如同时操作GPIO
// 加锁:
// - hwlocker:Lock('DO1', Self()) ->error
// 解锁:
// - hwlocker:UnLock('DO1', Self()) -> error
// 解锁:
// - hwlocker:Check('DO1', Self()) -> true|false

// 加锁
func LockHwResource(resName, luaUuid string) {

}

// 解锁
func UnLockHwResource(resName, luaUuid string) {

}

// 检查锁
func CheckHwResourceLock(resName, luaUuid string) bool {
	return false
}
