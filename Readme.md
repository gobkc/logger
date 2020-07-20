#logger

>wish list

- 浏览器
- ES        yes
- FILE      yes
- SYSLOG    yes
- MYSQL     no
- MONGO     no



logger接口

```
type Logger interface {
	Info(struPtr interface{})
	Warn(struPtr interface{})
	Error(struPtr interface{})
	Danger(struPtr interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}
```



> es demo

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
	
	


> file demo

	var logType = driver.FileLog{
		OutPut: "123.log",//如果这里的值指定空白，可以使用args的-log作为参数
	}
	logger.Set(logType)
	defer logType.Close()

> syslog demo

	var logType = driver.Syslog{
		Server:   "",
		Protocol: driver.SYS_LOC,
	}
	logger.Set(logType)

> mysql demo

	var logType = driver.Mysql{
		Server:   "",
	}
	logger.Set(logType)



使用说明：

第一步，选择使用的日志记录器类型，有elasticsearch、file、syslog等类型，将输出到指定对象

```
      logger.Set(logType)
```

第二步，获得日志器操作句柄或者logger

如果输出是elasticsearch, 则设置的是对应数据库，未做此步骤时，默认的数据库为defaultIndex；如果是文件则代表对应文件,

```
handle := GetLogHandle("test01")
```

第三步，日志输出

```
 m1 := Pet{
		ID:   id,
		Name: "logSelect",
		age:  rand.Intn(10000),
		Like: [3]string{"aaa", "bbb", "ddd"},
	}
	handle.Warn(&m1) //嵌套切片
handle.Println("hello")
```

