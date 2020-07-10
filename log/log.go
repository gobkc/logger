package log

import (
	"encoding/json"
	//"fmt"
	"log"
	"reflect"
	"time"
)

var logOprDefault logOperator

type OutPuter interface {
	OutPut(index string, p []byte) (n int, err error)
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
	return &logOprDefault
}

func createCombinationStru(v interface{}, add interface{}) (t reflect.Type) {
	fieldSlice := []reflect.StructField{}
	typ1 := reflect.TypeOf(v)
	typ2 := reflect.TypeOf(add)

	for i := 0; i < typ1.NumField(); i++ {
		fieldSlice = append(fieldSlice, typ1.Field(i))
	}
	for i := 0; i < typ2.NumField(); i++ {
		fieldSlice = append(fieldSlice, typ2.Field(i))
	}
	return reflect.StructOf(fieldSlice)
}

// copyDataToNewStru copy struct data to a struct
// newStInstance should be a pointer  value
func copyDataToNewStru(newStInstance interface{}, v interface{}) {
	rSrcValue := reflect.ValueOf(v)
	rSrcTyp := rSrcValue.Type()
	rInstance := reflect.ValueOf(newStInstance)

	if rInstance.Kind() != reflect.Ptr {
		return
	}
	rInstance = rInstance.Elem()

	var rDestField reflect.Value
	var rSrcField reflect.Value
	valueZero := reflect.Value{}
	// kind := rSrcValue.Kind()
	// if kind == reflect.Ptr {
	// 	kind = rSrcValue.Elem().Kind()
	// }
	switch rSrcValue.Kind() {
	case reflect.Struct, reflect.Ptr:
		for i := 0; i < rSrcValue.NumField(); i++ {
			rSrcField = rSrcValue.Field(i)
			//fmt.Println(name, rSrcField.Interface())
			rDestField = rInstance.FieldByName(rSrcTyp.Field(i).Name)
			if rDestField == valueZero {
				continue
			}
			if rSrcField.Type() != rDestField.Type() {
				//fmt.Println("type not same")
				continue
			}
			if rDestField.Kind() == reflect.Ptr {
				rSrcField = rSrcField.Elem()
			}
			if rDestField.Kind() == reflect.Struct ||
				rDestField.Kind() == reflect.Slice ||
				rDestField.Kind() == reflect.Array {
				if !rDestField.CanAddr() {
					continue
				}
				copyDataToNewStru(rDestField.Addr().Interface(), rSrcField.Interface())
				return
			}
			if !rDestField.CanSet() {
				//fmt.Println(rDestField, "cannot set")
				continue
			}

			copyReflectFieldValue(rDestField, rSrcField)
		}
	case reflect.Slice, reflect.Array:
		//rDestField = rInstance.FieldByName(rSrcTyp.Name())
		// if rDestField.Type() != rSrcField.Type() {
		// 	fmt.Println("slice type not same")
		// 	return
		// }
		//rDestField.SetCap(rSrcValue.Cap())
		//reflect.MakeSlice(rSrcField.Type(), rSrcValue.Len(), rSrcValue.Cap()
		//fmt.Println(rSrcValue.Len())
		rDestField = rInstance
		// // copy slice data
		for i := 0; i < rSrcValue.Len(); i++ {
			//fmt.Println(rSrcValue.Index(i).Interface())
			copyReflectFieldValue(rDestField.Index(i), rSrcValue.Index(i))
		}
	}
}

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
	l.output.OutPut(l.Index, data)
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
	logOprDefault.Error(v)
}

func Warn(v interface{}) {
	logOprDefault.Error(v)
}

func Println(v interface{}) {
	logOprDefault.Println(v)
}
