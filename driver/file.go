package driver

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileLog struct {
	OutPut     string
	FileHandle *os.File
}

//文件日志驱动
func (f *FileLog) InitFileLog() *os.File {
	f.ParseCmd()
	var file *os.File
	var err error
	//尝试创建全路径，有错误也不能返回
	var dir = filepath.Dir(f.OutPut)
	if err = os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		log.Println(dir, "已经存在")
	}
	if file, err = os.OpenFile(f.OutPut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Fatalln(err.Error())
	}
	f.FileHandle = file
	return file
}

//解析命令行   日志如果没有传递参数，会尝试从 cmd -log参数或则 log参数来获取
func (f *FileLog) ParseCmd() {
	var args = os.Args
	var output string
	for i, v := range args {
		var cond1 = strings.Contains(v, "log=")
		var cond2 = strings.Contains(v, "-log=")
		var cond3 = strings.Contains(v, "log")
		var cond4 = strings.Contains(v, "-log")
		if cond1 || cond2 {
			output = v[strings.Index(v, "=")+1:]
			break
		}
		if cond3 || cond4 {
			if len(args) >= i+1 {
				output = args[i+1]
			} else {
				output = ""
			}
			break
		}
	}
	if output != "" {
		f.OutPut = output
	}
}

//关闭文件句柄。通常在MAIN函数执行。否则还没写入日志就被关闭了
func (f *FileLog) Close() {
	if f.FileHandle != nil {
		f.FileHandle.Close()
	}
}
