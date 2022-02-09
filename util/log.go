package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)
			//fmt.Printf("filename=%v,function=%v line=%v\n", fileName, frame.Function, frame.Line)
			f := frame.Function
			ind := strings.LastIndex(f, "/")
			funcName := f[ind+1:]
			return "", fmt.Sprintf("%v(%v):%v", fileName, funcName, frame.Line)
		},
	})
}
