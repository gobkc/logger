package main

import (
	"log"
	"logger/logger"
	"logger/logger/driver"
)

func main() {
	//伪代码
	//var logType = driver.ElasticSearch{
	//	Host:     "127.0.0.1",
	//	Port:     9200,
	//	User:     "yunlifang",
	//	Password: "YlfEs2020",
	//	Index:    "test",
	//}
	var logType = driver.Syslog{
		Server:   "",
		Protocol: driver.SYS_LOC,
	}
	logger.Set(logType)
	log.Println(`{"msg":"abc"}`)
	log.Println(`ahahahahahah`)
}
