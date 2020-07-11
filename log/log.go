package log

import (
	"encoding/json"
	//"fmt"
	"log"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"
)

var logOprDefault logOperator

type OutPuter interface {
	OutPut(p []byte) (n int, err error)
	SetIndex(index string)
}

type EsHeader struct {
	Time  string `json:"time"`
	Level string `json:"level"`
}

type logOperator struct {
	Index  string
	output OutPuter
}

func SetOut(o OutPuter) {
	logOprDefault.output = o
}

func SetPrefix(index string) *logOperator {
	logOprDefault.Index = index
	if logOprDefault.output == nil {
		panic("didn't set ouput object")
	}
	logOprDefault.output.SetIndex(index)
	return &logOprDefault
}

func createCombinationStru(v interface{}, add interface{}) (t reflect.Type) {
	fieldSlice := []reflect.StructField{}
	typ1 := reflect.TypeOf(v)
	typ2 := reflect.TypeOf(add)

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
// src  : 源结构的接口
func copyDataToNewStru(dest interface{}, src interface{}) {
	rSrcValue := reflect.ValueOf(src)
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
	case reflect.String, reflect.Bool, reflect.Int,
		reflect.Float32, reflect.Float64:
		copyReflectFieldValue(rDestValue, rSrcValue)
	}
}

// copyReflectFieldValue标准变量值拷贝
func copyReflectFieldValue(rDestField reflect.Value, rSrcField reflect.Value) {
	if !rDestField.CanSet() {
		return
	}
	switch rDestField.Kind() {
	case reflect.String:
		rDestField.SetString(rSrcField.Interface().(string))
	case reflect.Bool:
		rDestField.SetBool(rSrcField.Interface().(bool))
	case reflect.Int:
		rDestField.SetInt(int64(rSrcField.Interface().(int)))
	case reflect.Float32:
		rDestField.SetFloat(float64(rSrcField.Interface().(float32)))
	case reflect.Float64:
		rDestField.SetFloat(float64(rSrcField.Interface().(float64)))
		// 其他还未支持
	}
}

func (l *logOperator) out(level string, v interface{}) {
	var esData = EsHeader{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Level: level,
	}

	// create a combination
	typ := createCombinationStru(v, esData)
	newStInstance := reflect.New(typ)
	copyDataToNewStru(newStInstance.Interface(), v)
	copyDataToNewStru(newStInstance.Interface(), esData)

	// marshal it
	data, _ := json.Marshal(newStInstance.Elem().Interface())
	l.output.OutPut(data)
}

// Error output info log
func (l *logOperator) Info(v interface{}) {
	l.out("info", v)
}

// Error output error log
func (l *logOperator) Error(v interface{}) {
	l.out("error", v)
}

// Danger output danger log
func (l *logOperator) Danger(v interface{}) {
	l.out("danger", v)
}

// Warn output warning log
func (l *logOperator) Warn(v interface{}) {
	l.out("warning", v)
}

func (l *logOperator) Println(v interface{}) {
	log.Println(v)
}

func Info(v interface{}) {
	logOprDefault.Info(v)
}

func Error(v interface{}) {
	logOprDefault.Error(v)
}

func Danger(v interface{}) {
	logOprDefault.Danger(v)
}

func Warn(v interface{}) {
	logOprDefault.Warn(v)
}

func Println(v interface{}) {
	logOprDefault.Println(v)
}
