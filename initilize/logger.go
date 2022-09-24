package initilize

import (
	"fmt"
	"hx/global"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLog() {
	c := global.Logger
	logger := logrus.New()
	logger.ReportCaller = true // 显示调用信息
	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	formatter.DisableQuote = true // 不转义换行符，为了保存错误堆栈到日志文件
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		return "", fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
	}
	logger.Formatter = formatter
	fileRotate := &lumberjack.Logger{
		Filename:   c.File,
		MaxBackups: 7,
	}
	writer := io.MultiWriter(os.Stdout, fileRotate)
	logger.SetOutput(writer)
	logLevel, _ := logrus.ParseLevel(c.LogLevel)
	logger.SetLevel(logLevel)
	logger.WithField("application", global.AppName)
	global.DL_LOGGER = logger
}
