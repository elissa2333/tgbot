package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// ToString 将任意类型转换为 string
func ToString(in interface{}) string {
	switch value := in.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case []byte:
		return string(value)
	case nil:
		return ""
	}

	return fmt.Sprintf("%v", in)
}

// ToInt 将任意类型转换为 int
func ToInt(in interface{}) int {
	switch value := in.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case nil:
		return 0
	}

	v, err := strconv.Atoi(ToString(in))
	if err != nil {
		return 0
	}

	return v
}

// StructToMap 将 struct 转化为 map
// 如果 struct field 拥有 Name 字段则自动 bind 为 struct 名 （json 库的坑）
func StructToMap(src interface{}) (map[string]interface{}, error) {
	refType := reflect.TypeOf(src)

	if refType.Kind() != reflect.Struct {
		if refType.Elem().Kind() == reflect.Struct {
			return StructToMap(reflect.ValueOf(src).Elem().Interface())
		}
		return nil, errors.New("input is not a struct")
	}

	temp, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(temp, &result)
	return result, err
}
