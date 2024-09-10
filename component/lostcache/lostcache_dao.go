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

package lostcache

/**
 * Save
 *
 */
func SaveLostCacheData(data CacheDataDto) error {
	__Sqlite.db.Table("cache_data").Create(&CacheData{
		TargetId: data.TargetId,
		Data:     data.Data,
	})
	return __Sqlite.db.Error
}

/**
 * Get
 *
 */
func GetLostCacheData(uuid string) ([]CacheDataDto, error) {
	dataDto := []CacheDataDto{}
	__Sqlite.db.Model(&CacheData{}).Where("target_id=?", uuid).Find(&dataDto)
	return dataDto, __Sqlite.db.Error

}

/**
 * Delete Lost Cache Data
 *
 */
func DeleteLostCacheData(dbId uint) {
	__Sqlite.db.Model(&CacheData{}).Where("id=?", dbId).Delete(&CacheData{})
}

/**
 * Clear
 *
 */
func ClearLostCacheData(uuid string) {
	__Sqlite.db.Model(&CacheData{}).Where("target_id=?", uuid).Delete(&CacheData{})
}
