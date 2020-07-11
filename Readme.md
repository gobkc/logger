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
		logger.Set(logType)
	
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
	
		log.SetPrefix("test01").Info(m)
	
		id++
		m1 := Pet{
			ID:   id,
			Name: "wangwang",
			age:  rand.Intn(10000),
			Like: [3]string{"aaa", "bbb", "ddd"},
		}
		log.SetPrefix("test01").Info(m1)
		id++
		m2 := Messages{
			ID: id,
			Msgs: []Message{
				m, m, m,
			},
		}
		log.Danger(m2)

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

