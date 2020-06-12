package logger

import (
	"fmt"
	"log"
	"logger/logger/driver"
)

func Set(logType interface{}) {
	t := fmt.Sprintf("%T", logType)
	switch t {
	case "driver.ElasticSearch":
		var to = logType.(driver.ElasticSearch)
		log.SetFlags(0)
		log.SetOutput(&to)
	case "driver.Syslog":
		var to = logType.(driver.Syslog)
		log.SetOutput(&to)
	}
}
