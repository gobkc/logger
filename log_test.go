package logger

import (
	"github.com/gobkc/logger/driver"
	"math/rand"
	"testing"
	"time"
)

func TestDemo01(t *testing.T) {
	rand.Seed(time.Now().Unix())
	type BookInfo struct {
		ID     int    `json:"id"`
		Author string `json:"author"`
		Sales  int    `json:"sales"`
	}

	type Message struct {
		ID   int      `json:"id"`
		Name string   `json:"name1"`
		Age  int      `json:"age"`
		Book BookInfo `json:"book"`
	}

	type Messages struct {
		ID   int       `json:"id"`
		Name string    `json:"name1"`
		Msgs []Message `json:"memgs"`
	}

	type Pet struct {
		ID   int       `json:"id"`
		Name string    `json:"petName"`
		age  int       `json:"petAge"`
		Like [3]string `json:"like"`
	}
	var logType = driver.ElasticSearch{
		Common: &driver.ElasticSearchCommon{
			Host:     "89zx.com",
			Port:     9200,
			User:     "yunlifang",
			Password: "YlfEs2020",
		},
	}

	id := int(time.Now().Unix())

	m := Message{
		ID:   id,
		Name: "logSelect",
		Age:  rand.Intn(10000),
		Book: BookInfo{
			Author: "name 2",
			Sales:  rand.Intn(10000),
		},
	}

	Set(logType)
	//嵌套结构体
	handle := GetLogHandle("test01")
	handle.Info(&m)
	id++
	m1 := Pet{
		ID:   id,
		Name: "logSelect",
		age:  rand.Intn(10000),
		Like: [3]string{"aaa", "bbb", "ddd"},
	}
	handle.Warn(&m1) //嵌套切片
	id++
	m2 := Messages{
		ID:   id,
		Name: "logselect",
		Msgs: []Message{
			m, m, m,
		},
	}

	handle.Danger(&m2) //嵌套切片结构
	id++
	m3 := Pet{
		ID:   id,
		Name: "logSelect",
		age:  rand.Intn(10000),
		Like: [3]string{"aaa", "bbb", "ddd"},
	}
	handle.Error(&m3)
	handle.Print("you are login in")
}
