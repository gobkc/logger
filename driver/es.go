package driver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const defaultIndex = "defaultIndex"

// ElasticSearch log driver
type ElasticSearch struct {
	Common *ElasticSearchCommon
	Index  string
}

// ElasticSearchCommon 通用信息
type ElasticSearchCommon struct {
	Host     string
	Port     int
	User     string //basic auth 如果没有账号密码请留空
	Password string //basic auth
	Https    bool
}

// EsDefault 当输入日志不是一个JSON的默认处理方式
type EsDefault struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

// EsHeader 日志消息头
type EsHeader struct {
	Time  string `json:"time"`
	Level string `json:"level"`
}

func (e *ElasticSearch) out(level string, v interface{}) {
	var esData = EsHeader{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Level: level,
	}

	// create a combination
	newStInstance := createCombinationStru(v, &esData)
	copyDataToNewStru(newStInstance.Interface(), v)
	copyDataToNewStru(newStInstance.Interface(), &esData)

	// marshal it
	data, err := json.Marshal(newStInstance.Elem().Interface())
	if err != nil {
		return
	}

	e.MultiWrite(data)
}

// Info output info log
func (e *ElasticSearch) Info(v interface{}) {
	e.out("info", v)
}

// Error output error log
func (e *ElasticSearch) Error(v interface{}) {
	e.out("error", v)
}

// Danger output danger log
func (e *ElasticSearch) Danger(v interface{}) {
	e.out("danger", v)
}

// Warn output warning log
func (e *ElasticSearch) Warn(v interface{}) {
	e.out("warning", v)
}

// Println 取得打印信息作为日志message输出
func (e *ElasticSearch) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	e.MultiWrite(data)
}

// Print 取得打印信息作为日志message输出
func (e *ElasticSearch) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	e.MultiWrite(data)
}

// Printf 取得打印信息作为日志message输出
func (e *ElasticSearch) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	e.MultiWrite(data)
}

// MultiWrite 向日志记录器输出
func (e *ElasticSearch) MultiWrite(p []byte) (n int, err error) {
	// 将数据输出到多个输出目的地,采用io.MultiWriter存在问题如果写入一个失败，则不能继续写
	str := string(p)
	fmt.Fprint(os.Stdout, str)
	fmt.Fprint(e, str)
	return len(p), err
}

func packEsMessage(p []byte) (data []byte, err error) {
	var newParm = EsDefault{
		Message: string(p),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Level:   "info",
	}

	data, err = json.Marshal(newParm)
	return
}

// Write 根据数据和index输出到elasticsearch
func (e *ElasticSearch) Write(p []byte) (n int, err error) {
	var eType = fmt.Sprintf("%s:%s", e.Common.User, e.Common.Password)
	var esEncode = base64.StdEncoding.EncodeToString([]byte(eType))
	var auth = fmt.Sprintf("Basic %s", esEncode)
	var param = bytes.NewReader(p)
	var actMethod = "http"
	if e.Common.Https {
		actMethod = "https"
	}
	// fmt.Println("write:", string(p))
	var api = fmt.Sprintf("%s://%s:%d/%s/_doc", actMethod, e.Common.Host, e.Common.Port, e.Index)
	var request = new(http.Request)
	request, err = http.NewRequest("POST", api, param)
	if err != nil {
		return n, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", auth)
	var response = new(http.Response)
	response, err = http.DefaultClient.Do(request)
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	if response == nil || response.Body == nil {
		return 0, errors.New("elasticsearch host error")
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return n, err
	}
	var noError = err
	return len(p), noError
}
