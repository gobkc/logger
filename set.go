package logger

import (
	"fmt"
	"github.com/gobkc/logger/driver"
	mylog "github.com/gobkc/logger/log"
	"log"
)

func Set(logType interface{}) {
	t := fmt.Sprintf("%T", logType)
	switch t {
	case "driver.ElasticSearch":
		var to = logType.(driver.ElasticSearch)
		log.SetFlags(0)
		log.SetOutput(&to)
		mylog.SetOut(&to)
	case "driver.Syslog":
		var to = logType.(driver.Syslog)
		log.SetFlags(0)
		log.SetOutput(&to)
	case "driver.FileLog":
		var to = logType.(driver.FileLog)
		log.SetFlags(0)
		var file = to.InitFileLog()
		log.SetOutput(file)
	}
}
