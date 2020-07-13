#logger

>wish list

- 浏览器
- ES        yes
- FILE      yes
- SYSLOG    yes
- MYSQL     no
- MONGO     no

> es demo

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
			Msgs []Message `json:"memgs"`
		}
	
		type Pet struct {
			ID   int       `json:"id"`
			Name string    `json:"petName"`
			age  int       `json:"petAge"`
			Like [3]string `json:"like"`
		}
		var logType = driver.ElasticSearch{
			Host:     "hello", //"89zx.com",
			Port:     9200,
			User:     "zou", //"yunlifang",
			Password: "YlfEs2020",
		}
	
		id := int(time.Now().Unix())
	
		m := Message{
			ID:   id,
			Name: "li",
			Age:  rand.Intn(10000),
			Book: BookInfo{
				Author: "name 2",
				Sales:  rand.Intn(10000),
			},
		}
	
		logger.Set(logType)
		logger.Println("hello")
		logger.SetPrefix("test01").Info(&m)
	
		id++
		m1 := Pet{
			ID:   id,
			Name: "wangwang",
			age:  rand.Intn(10000),
			Like: [3]string{"aaa", "bbb", "ddd"},
		}
		logger.SetPrefix("test01").Info(&m1)
		id++
		m2 := Messages{
			ID: id,
			Msgs: []Message{
				m, m, m,
			},
		}
		logger.Danger(&m2)
		s := "hello"
		logger.Println(s)
	

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

第二步，设置prefix

如果输出是elasticsearch, 则设置的是对应数据库，未做此步骤时，默认的数据库为defaultIndex；如果是文件则代表对应文件,

```
logger.SetPrefix("test01")
```

第三步，日志输出

```
#方式1，通过设置prefix返回的引用调用日志记录方法,此时传入的信息格式为结构体的指针
logger.SetPrefix("test01").Info(&m)
#方式2.调用logger包函数print,传入数据为字符串，作为单信息message输出
logger.Println("hello")
logger.Print("hi")
```

