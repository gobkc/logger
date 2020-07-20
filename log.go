package logger

import (
	"github.com/gobkc/logger/driver"
)

// logType 记录选择的日志记录器类型
var logType interface{}

// Logger 日志器对外接口
type Logger interface {
	Info(v interface{})
	Warn(v interface{})
	Error(v interface{})
	Danger(v interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}

// setLogType 获得日志驱动的引用(driver指针类型)
func setLogType(o interface{}) {
	logType = o
}

// GetLogHandle 设置写的目标，elasticSearch则指数据库
func GetLogHandle(index string) Logger {
	var logger Logger
	if logType == nil {
		panic("you should set logType Before")
	}
	switch logType.(type) {
	case *driver.ElasticSearch:
		elastic := driver.ElasticSearch{Index: index}
		//common := getEsCommonInstance(logType.(*Elastic))
		esInfo := logType.(*driver.ElasticSearch)
		elastic.Common = (*esInfo).Common
		logger = &elastic
	case *driver.Mysql:
		//todo
		logger = nil
	}
	return logger
}
