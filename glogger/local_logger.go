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

package glogger

import (
	"os"
)

// LogWriter 定义了一个本地日志写入器
type LogWriter struct {
	file *os.File
}

// NewLogWriter 创建一个新的 LogWriter 实例
func NewLogWriter(filepath string) *LogWriter {
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		GLogger.Fatalf("Fail to open log file %s: %v", filepath, err)
		os.Exit(1)
	}
	return &LogWriter{file: logFile}
}

// Write 将字节切片写入日志文件
func (lw *LogWriter) Write(b []byte) (n int, err error) {
	return lw.file.Write(b)
}

// Close 关闭日志文件
func (lw *LogWriter) Close() error {
	if lw.file != nil {
		return lw.file.Close()
	}
	return nil
}
