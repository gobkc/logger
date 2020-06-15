#logger

>wish list

- 浏览器
- ES        yes
- FILE      yes
- SYSLOG    yes
- MYSQL     no
- MONGO     no

> es demo

	var logType = driver.ElasticSearch{
		Host:     "127.0.0.1",
		Port:     9200,
		User:     "yunlifang",
		Password: "YlfEs2020",
		Index:    "test",
	}
    logger.Set(logType)

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

