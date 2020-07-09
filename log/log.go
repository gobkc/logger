package log

import (
	"log"
)

type EsDefault struct {
	Time    string      `json:"time"`
	Message interface{} `json:"message"`
	Level   string      `json:"level"`
}

func Error(v interface{}) {
	var esData = EsDefault{
		Message: v,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "error",
	}
	log.Println(esData)
}

func Info(v interface{}) {
	var esData = EsDefault{
		Message: v,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "info",
	}
	log.Println(esData)
}

func Danger(v interface{}) {
	var esData = EsDefault{
		Message: v,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "danger",
	}
	log.Println(esData)
}

func Warn(v interface{}) {
	var esData = EsDefault{
		Message: v,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "warning",
	}
	log.Println(esData)
}
