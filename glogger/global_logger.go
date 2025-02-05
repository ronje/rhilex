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

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 定义了日志配置结构体
type LogConfig struct {
	AppID         string
	LogLevel      string
	EnableConsole bool
	DebugMode     bool
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	LogCompress   bool
}

// Logrus 是 logrus 的全局实例
var Logrus *logrus.Logger

// GLogger 是 logrus 的全局 Entry 实例
var GLogger *logrus.Entry

// StartGLogger 启动全局日志记录器
func StartGLogger(config LogConfig) {
	Logrus = logrus.New()
	GLogger = Logrus.WithField("appId", config.AppID)
	Logrus.Formatter = &logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   "2006-01-02T15:04:05.999999999Z07:00",
	}
	if config.DebugMode {
		Logrus.SetReportCaller(true)
	}
	if config.EnableConsole {
		Logrus.SetOutput(os.Stdout)
	} else {
		Logrus.SetOutput(&lumberjack.Logger{
			Filename:   "rhilex-running-log.txt",
			MaxSize:    config.LogMaxSize,
			MaxBackups: config.LogMaxBackups,
			MaxAge:     config.LogMaxAge,
			Compress:   config.LogCompress,
		})
	}
	setLogLevel(config.LogLevel)
}

// setLogLevel 设置日志级别
func setLogLevel(logLevel string) {
	levelMap := map[string]logrus.Level{
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"all":   logrus.TraceLevel,
		"trace": logrus.TraceLevel,
	}
	if level, ok := levelMap[logLevel]; ok {
		Logrus.SetLevel(level)
	} else {
		Logrus.SetLevel(logrus.InfoLevel)
	}
}

// Info 记录 INFO 级别的日志
func Info(args ...interface{}) {
	GLogger.Info(args...)
}

// Infof 记录 INFO 级别的格式化日志
func Infof(format string, args ...interface{}) {
	GLogger.Infof(format, args...)
}

// Error 记录 ERROR 级别的日志
func Error(args ...interface{}) {
	GLogger.Error(args...)
}

// Errorf 记录 ERROR 级别的格式化日志
func Errorf(format string, args ...interface{}) {
	GLogger.Errorf(format, args...)
}

// Debug 记录 DEBUG 级别的日志
func Debug(args ...interface{}) {
	GLogger.Debug(args...)
}

// Debugf 记录 DEBUG 级别的格式化日志
func Debugf(format string, args ...interface{}) {
	GLogger.Debugf(format, args...)
}

// Close 关闭日志记录器
func Close() error {
	return GLogger.Writer().Close()
}
