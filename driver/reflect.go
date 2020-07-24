package driver

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

func createCombinationStru(v interface{}, add interface{}) (t reflect.Value) {
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
	newtyp := reflect.StructOf(fieldSlice)
	return reflect.New(newtyp)
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
	//fmt.Println("copyReflectFieldValue", rSrcField.Interface())
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
