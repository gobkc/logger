#logger

>wish list

- 浏览器
- ES        yes
- FILE      yes
- SYSLOG    yes
- MYSQL     no
- MONGO     no

> es demo

		rand.Seed(time.Now().Unix())
		type BookInfo struct {
			Author string `json:"author"`
			Sales  int    `json:"sales"`
		}
	
		type Message struct {
			Name string   `json:"name1"`
			Age  int      `json:"age"`
			Book BookInfo `json:"book"`
		}
	
		type Pet struct {
			Name string    `json:"petName"`
			Age  int       `json:"petAge"`
			Like [3]string `json:"like"`
		}
		var logType = driver.ElasticSearch{
			Host:     "89zx.com",
			Port:     9200,
			User:     "yunlifang",
			Password: "YlfEs2020",
			Index:    "test",
		}
		logger.Set(logType)
	
		m := Message{
			Name: "li",
			Age:  rand.Intn(10000),
			Book: BookInfo{
				Author: "name 2",
				Sales:  rand.Intn(10000),
			},
		}
	
		log.SetPrefix("test").Info(m)
		// log.SetPrefix("test").Warn(m)
		// log.SetPrefix("test").Error(m)
		// log.SetPrefix("test").Danger(m)
		// log.Warn(m)
		// log.Info(m)
		m1 := Pet{
			Name: "wangwang",
			Age:  rand.Intn(10000),
			Like: [3]string{"aaa", "bbb", "ddd"},
		}
		log.SetPrefix("test").Info(m1)

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

