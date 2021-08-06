package main

import "reflect"

// ExtractStrings 从任意类型提取出所有 string 类型的字段
func ExtractStrings(x interface{}, fn func(string)) {
	val := reflect.ValueOf(x)

	// use the actual type pointed to by the pointer
	// 如果是指针类型, 调用 Elem() 得到实际类型
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 提取值
	extractValue := func(value reflect.Value) {
		// private field cannot call Interface()
		// 私有字段不能转 Interface
		if value.CanInterface() {
			ExtractStrings(value.Interface(), fn)
		}
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			extractValue(val.Field(i))
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			extractValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			extractValue(val.MapIndex(key))
		}
	}
}
