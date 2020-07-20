package logger

import (
	"fmt"
	"github.com/gobkc/logger/driver"
	"log"
)

// Set select a log driver, and set loginfo
func Set(logType interface{}) {
	t := fmt.Sprintf("%T", logType)
	switch t {
	case "driver.ElasticSearch":
		var to = logType.(driver.ElasticSearch)
		setLogType(&to)
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
