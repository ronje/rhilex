package glogger

import (
	"os"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
*
* 配置全局logging记录器
*
 */

var Logrus *logrus.Logger
var GLogger *logrus.Entry

func StartGLogger(appId string, LogLevel string, EnableConsole bool,
	DebugMode bool, LogMaxSize, LogMaxBackups, LogMaxAge int, LogCompress bool) {
	Logrus = logrus.New()
	GLogger = Logrus.WithField("appId", appId)
	Logrus.Formatter = &logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   "2006-01-02T15:04:05.999999999Z07:00",
	}
	if DebugMode {
		Logrus.SetReportCaller(true)
	}
	if EnableConsole {
		Logrus.SetOutput(os.Stdout)
	} else {
		Logrus.SetOutput(&lumberjack.Logger{
			Filename:   "rhilex-running-log.txt",
			MaxSize:    LogMaxSize,    // 超过10Mb备份
			MaxBackups: LogMaxBackups, // 最多备份3次
			MaxAge:     LogMaxAge,     // 最大保留天数
			Compress:   LogCompress,   // 压缩备份
		})
	}

	setLogLevel(LogLevel)
}
func setLogLevel(LogLevel string) {
	switch LogLevel {
	case "fatal":
		Logrus.SetLevel(logrus.FatalLevel)
	case "error":
		Logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		Logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		Logrus.SetLevel(logrus.DebugLevel)
	case "info":
		Logrus.SetLevel(logrus.InfoLevel)
	case "all", "trace":
		Logrus.SetLevel(logrus.TraceLevel)
	}

}
func Info(args ...interface{}) {
	GLogger.Info(args...)
}
func Infof(format string, args ...interface{}) {
	GLogger.Infof(format, args...)
}
func Error(args ...interface{}) {
	GLogger.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	GLogger.Errorf(format, args...)
}
func Debug(args ...interface{}) {
	GLogger.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	GLogger.Debugf(format, args...)
}

/*
*
* 关闭日志记录器
*
 */
func Close() error {
	return GLogger.Writer().Close()
}
