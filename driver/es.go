package driver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultIndex = "defaultIndex"

// ElasticSearch log driver
type ElasticSearch struct {
	Host     string
	Port     int
	User     string //basic auth 如果没有账号密码请留空
	Password string //basic auth
	index    string
	Https    bool
}

// //当输入日志不是一个JSON的默认处理方式
// type EsDefault struct {
// 	Time    string `json:"time"`
// 	Message string `json:"message"`
// 	Level   string `json:"level"`
// }

func (e *ElasticSearch) SetIndex(index string) {
	e.index = index
}

//WriteMessage 将消息p组装到完整结构中进行输出
// func (e *ElasticSearch) WriteMessage(p []byte) (n int, err error) {
// 	var parseLog interface{}
// 	//先用json尝试解析，如果解析不了，则使用自创建结构体来解析
// 	if err = json.Unmarshal(p, &parseLog); err != nil {
// 		var newParm = EsDefault{
// 			Message: string(p),
// 			Time:    time.Now().Format("2006-01-02 15:04:05"),
// 			Level:   "info",
// 		}
// 		p, err = json.Marshal(newParm)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	if e.index == "" {
// 		e.index = defaultIndex
// 	}
// 	return e.Write(p)
// }

// Write 根据数据和index输出到elasticsearch
func (e *ElasticSearch) Write(p []byte) (n int, err error) {
	var eType = fmt.Sprintf("%s:%s", e.User, e.Password)
	var esEncode = base64.StdEncoding.EncodeToString([]byte(eType))
	var auth = fmt.Sprintf("Basic %s", esEncode)
	var param = bytes.NewReader(p)
	var actMethod = "http"
	if e.Https {
		actMethod = "https"
	}

	var api = fmt.Sprintf("%s://%s:%d/%s/_doc", actMethod, e.Host, e.Port, e.index)
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
