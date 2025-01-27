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
