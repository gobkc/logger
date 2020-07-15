package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"
)

var logOprDefault LogOperator

// OutPuter 具体驱动的接口
type OutPuter interface {
	Write(p []byte) (n int, err error)
	SetIndex(index string)
}

// EsHeader 日志消息头
type EsHeader struct {
	Time  string `json:"time"`
	Level string `json:"level"`
}

// EsDefault 当输入日志不是一个JSON的默认处理方式
type EsDefault struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

// LogOperator 日志操作器，对外提供接口
type LogOperator struct {
	Index    string
	outputer OutPuter
}

// SetOut 获得日志驱动的引用
func SetOut(o OutPuter) {
	logOprDefault.outputer = o
}

// SetPrefix 设置写的目标，elasticSearch则指数据库
func SetPrefix(index string) *LogOperator {
	logOprDefault.Index = index
	if logOprDefault.outputer == nil {
		panic("didn't  Set driver")
	}
	logOprDefault.outputer.SetIndex(index)
	return &logOprDefault
}

func createCombinationStru(v interface{}, add interface{}) (t reflect.Type) {
	fieldSlice := []reflect.StructField{}
	typ1 := reflect.TypeOf(v).Elem()
	typ2 := reflect.TypeOf(add).Elem()

	for i := 0; i < typ1.NumField(); i++ {
		// discard lowercase field
		r, _ := utf8.DecodeRuneInString(typ1.Field(i).Name)
		if unicode.IsLower(r) {
			//fmt.Printf("%v, %v, %c\n", typ1.Field(i).Name, "islower:", r)
			continue
		}
		fieldSlice = append(fieldSlice, typ1.Field(i))
	}
	for i := 0; i < typ2.NumField(); i++ {
		// discard lowercase field
		r, _ := utf8.DecodeRuneInString(typ2.Field(i).Name)
		if unicode.IsLower(r) {
			//fmt.Println(typ2.Field(i).Name, "islower:", r)
			continue
		}
		fieldSlice = append(fieldSlice, typ2.Field(i))
	}
	return reflect.StructOf(fieldSlice)
}

// copyDataToNewStru copy struct data to a struct
// dest : 目标结构的指针的接口
// src  : 源结构指针的接口
func copyDataToNewStru(dest interface{}, src interface{}) {
	rSrcValue := reflect.ValueOf(src).Elem()
	rSrcTyp := rSrcValue.Type()
	rDestValue := reflect.ValueOf(dest)

	if rDestValue.Kind() != reflect.Ptr {
		return
	}
	rDestValue = rDestValue.Elem()
	if rSrcTyp.Kind() == reflect.Ptr {
		rSrcValue = rSrcValue.Elem()
	}

	var rDestField reflect.Value
	var rSrcField reflect.Value
	valueZero := reflect.Value{}

	for i := 0; i < rSrcValue.NumField(); i++ {
		rSrcField = rSrcValue.Field(i)
		rDestField = rDestValue.FieldByName(rSrcTyp.Field(i).Name)
		// not found
		if rDestField == valueZero {
			continue
		}
		if rSrcField.Type() != rDestField.Type() {
			continue
		}

		if !rDestField.CanSet() {
			//fmt.Println(rDestField, "cannot set")
			continue
		}
		copyEquivalentElement(rDestField, rSrcField)
	}

}

// copyEquivalentElement 支持结构、切片、以及他们的嵌套、标量拷贝
func copyEquivalentElement(rDestValue reflect.Value, rSrcValue reflect.Value) {
	// 过滤不同类型的copy,未产生panic
	if rDestValue.Type() != rSrcValue.Type() {
		return
	}

	var rDestField reflect.Value
	var rSrcField reflect.Value

	switch rSrcValue.Kind() {
	case reflect.Struct:
		for i := 0; i < rSrcValue.NumField(); i++ {
			rSrcField = rSrcValue.Field(i)
			rDestField = rDestValue.Field(i)
			// struct field not support pointer now
			if rDestField.Kind() == reflect.Struct ||
				rDestField.Kind() == reflect.Slice ||
				rDestField.Kind() == reflect.Array {
				copyEquivalentElement(rDestField, rSrcField)
				return
			}
			if !rDestField.CanSet() {
				continue
			}
			copyEquivalentElement(rDestField, rSrcField)
		}
	case reflect.Slice, reflect.Array:
		if rSrcValue.Kind() == reflect.Slice {
			newVal := reflect.MakeSlice(rDestValue.Type(), rSrcValue.Cap(), rSrcValue.Cap())
			rDestValue.Set(newVal)
		}
		// // copy slice data
		for i := 0; i < rSrcValue.Len(); i++ {
			copyEquivalentElement(rDestValue.Index(i), rSrcValue.Index(i))
		}
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		copyReflectFieldValue(rDestValue, rSrcValue)
	}
}

// copyReflectFieldValue标准变量值拷贝
func copyReflectFieldValue(rDestField reflect.Value, rSrcField reflect.Value) {
	if !rDestField.CanSet() {
		return
	}
	fmt.Println("copyReflectFieldValue", rSrcField.Interface())
	srcValue := rSrcField.Interface()
	switch v := srcValue.(type) {
	case string:
		rDestField.SetString(v)
	case bool:
		rDestField.SetBool(v)
	case int, int8, int16, int32, int64:
		if itsValue, ok := v.(int); ok {
			rDestField.SetInt(int64(itsValue))
		} else if itsValue, ok := v.(int8); ok {
			rDestField.SetInt(int64(itsValue))
		} else if itsValue, ok := v.(int16); ok {
			rDestField.SetInt(int64(itsValue))
		} else if itsValue, ok := v.(int32); ok {
			rDestField.SetInt(int64(itsValue))
		} else if itsValue, ok := v.(int64); ok {
			rDestField.SetInt(int64(itsValue))
		}
	case uint, uint8, uint16, uint32, uint64, uintptr:
		if itsValue, ok := v.(uint); ok {
			rDestField.SetUint(uint64(itsValue))
		} else if itsValue, ok := v.(uint8); ok {
			rDestField.SetUint(uint64(itsValue))
		} else if itsValue, ok := v.(uint16); ok {
			rDestField.SetUint(uint64(itsValue))
		} else if itsValue, ok := v.(uint32); ok {
			rDestField.SetUint(uint64(itsValue))
		} else if itsValue, ok := v.(uint64); ok {
			rDestField.SetUint(uint64(itsValue))
		} else if itsValue, ok := v.(uintptr); ok {
			rDestField.SetUint(uint64(itsValue))
		}
	case float32, float64:
		if itsValue, ok := v.(uintptr); ok {
			rDestField.SetFloat(float64(itsValue))
		} else if itsValue, ok := v.(uintptr); ok {
			rDestField.SetFloat(float64(itsValue))
		}
		// 其他还未支持
	}
}

func (l *LogOperator) out(level string, v interface{}) {
	rTyp := reflect.TypeOf(v)
	if rTyp.Kind() != reflect.Ptr ||
		rTyp.Elem().Kind() != reflect.Struct {
		panic("log data type is not support")
		return
	}
	var esData = EsHeader{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Level: level,
	}

	// create a combination
	typ := createCombinationStru(v, &esData)
	newStInstance := reflect.New(typ)
	copyDataToNewStru(newStInstance.Interface(), v)
	copyDataToNewStru(newStInstance.Interface(), &esData)

	// marshal it
	data, err := json.Marshal(newStInstance.Elem().Interface())
	if err != nil {
		return
	}

	l.Write(data)
}

// Error output info log
func (l *LogOperator) Info(v interface{}) {
	l.out("info", v)
}

// Error output error log
func (l *LogOperator) Error(v interface{}) {
	l.out("error", v)
}

// Danger output danger log
func (l *LogOperator) Danger(v interface{}) {
	l.out("danger", v)
}

// Warn output warning log
func (l *LogOperator) Warn(v interface{}) {
	l.out("warning", v)
}

// Println 取得打印信息作为日志message输出
func (l *LogOperator) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	l.Write(data)
}

// Print 取得打印信息作为日志message输出
func (l *LogOperator) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	l.Write(data)
}

// Printf 取得打印信息作为日志message输出
func (l *LogOperator) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	data, err := packEsMessage([]byte(message))
	if err != nil {
		return
	}
	l.Write(data)
}

// Write 向日志记录器输出
func (l *LogOperator) Write(p []byte) (n int, err error) {
	// 将数据输出到多个输出目的地
	reader := bytes.NewReader(p)
	mutiWriter := io.MultiWriter(l.outputer, os.Stdout)
	num, err := io.Copy(mutiWriter, reader)
	return int(num), err
}

// Info logOperator.Info 的封装
func Info(v interface{}) {
	logOprDefault.Info(v)
}

// Error logOperator.Error 的封装
func Error(v interface{}) {
	logOprDefault.Error(v)
}

// Danger logOperator.Danger 的封装
func Danger(v interface{}) {
	logOprDefault.Danger(v)
}

// Warn logOperator.Warn 的封装
func Warn(v interface{}) {
	logOprDefault.Warn(v)
}

// Println logOperator.Println 的封装
func Println(v ...interface{}) {
	if logOprDefault.outputer == nil {
		panic("didn't Set output driver")
	}
	logOprDefault.Println(v...)
}

// Print logOperator.Print 的封装
func Print(v ...interface{}) {
	if logOprDefault.outputer == nil {
		panic("didn't Set output driver")
	}
	logOprDefault.Print(v...)
}

// Printf logOperator.Printf 的封装
func Printf(format string, v ...interface{}) {
	if logOprDefault.outputer == nil {
		panic("didn't Set output driver")
	}
	logOprDefault.Printf(format, v...)
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
