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

package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCsvToMap reads a CSV file and returns a slice of maps with header as keys
func ReadCsvToMap(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("insufficient data in CSV")
	}

	csvHeader := records[0]
	dataRows := records[1:]

	var result []map[string]string

	for _, row := range dataRows {
		if len(row) != len(csvHeader) {
			return nil, fmt.Errorf("row length mismatch")
		}

		rowData := make(map[string]string)
		for i, columnName := range csvHeader {
			rowData[columnName] = row[i]
		}
		result = append(result, rowData)
	}

	return result, nil
}

func TestCSV() {
	filePath := "registers.csv"
	mapData, err := ReadCsvToMap(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("ReadCsvToMap Result:")
	for _, row := range mapData {
		fmt.Println(row)
	}
}
