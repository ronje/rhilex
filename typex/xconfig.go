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

package typex

// Global config
type ExtLib struct {
	Value []string `ini:"value,,allowshadow"`
}
type Secret struct {
	Value []string `ini:"value,,allowshadow"`
}
type RhilexConfig struct {
	IniPath               string   `json:"-"`
	AppId                 string   `ini:"app_id" json:"appId"`
	MaxQueueSize          int      `ini:"max_queue_size" json:"maxQueueSize"`
	SourceRestartInterval int      `ini:"resource_restart_interval" json:"sourceRestartInterval"`
	GomaxProcs            int      `ini:"gomax_procs" json:"gomaxProcs"`
	EnablePProf           bool     `ini:"enable_pprof" json:"enablePProf"`
	EnableConsole         bool     `ini:"enable_console" json:"enableConsole"`
	DebugMode             bool     `ini:"debug_mode" json:"appDebugMode"`
	LogLevel              string   `ini:"log_level" json:"logLevel"`
	LogMaxSize            int      `ini:"log_max_size" json:"logMaxSize"`
	LogMaxBackups         int      `ini:"log_max_backups" json:"logMaxBackups"`
	LogMaxAge             int      `ini:"log_max_age" json:"logMaxAge"`
	LogCompress           bool     `ini:"log_compress" json:"logCompress"`
	MaxKvStoreSize        int      `ini:"max_kv_store_size" json:"maxKvStoreSize"`
	ExtLibs               []string `ini:"ext_libs,,allowshadow" json:"extLibs"`
	DataSchemaSecret      []string `ini:"dataschema_secrets,,allowshadow" json:"dataSchemaSecret"`
}
