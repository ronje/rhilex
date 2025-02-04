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
	"testing"

	"github.com/hootrhino/rhilex/datacenter"
)

// go test -timeout 30s -run ^Test_gen_schema_ddl github.com/hootrhino/rhilex/test -v -count=1
func Test_gen_schema_ddl(t *testing.T) {
	schema := datacenter.SchemaDDL{
		SchemaUUID: "ABCDEF",
		DDLColumns: []datacenter.DDLColumn{
			{Name: "id", Type: "int", Description: "Pk"},
			{Name: "ts", Type: "int", Description: "Timestamp"},
			{Name: "schema_id", Type: "int", Description: "schema id"},
			{Name: "temp", Type: "float", Description: "temp"},
			{Name: "humi", Type: "float", Description: "humi"},
			{Name: "status", Type: "bool", Description: "humi"},
			{Name: "msg", Type: "string", Description: "humi"},
		},
	}
	createTableSQL, err1 := datacenter.GenerateSQLiteCreateTableDDL(schema)
	if err1 != nil {
		t.Log("Error:", err1)
		return
	}
	t.Log(createTableSQL)
}
