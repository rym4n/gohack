package gohack

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

//Round 保留小数点后几位
func Round(num float64, precise int) float64 {
	preciseStr := fmt.Sprintf("%d", precise)
	tpl := "%0." + preciseStr + "f"

	s := fmt.Sprintf(tpl, num)
	bitSize := 64
	ret, _ := strconv.ParseFloat(s, bitSize)
	return ret
}

// InArray 判断某个值是否在数组中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// StrVal 把其他类型转为string
func StrVal(val interface{}) string {
	var strTmp string
	valueType := reflect.TypeOf(val).Kind()
	valueValue := reflect.ValueOf(val)
	switch valueType {
	case reflect.Map, reflect.Slice, reflect.Array:
		strTmp = "can't convert to string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		strTmp = fmt.Sprintf("%d", valueValue)
	case reflect.Float32, reflect.Float64:
		strTmp = fmt.Sprintf("%f", valueValue)
	case reflect.String:
		strTmp = fmt.Sprintf("%s", valueValue)
	default:
		strTmp = ""
	}
	return strTmp
}

// IntVal 把其他类型转为int
func IntVal(val interface{}) int {
	var intTmp int
	valueType := reflect.TypeOf(val).Kind()
	switch valueType {
	case reflect.Slice, reflect.Array, reflect.Map:
		intTmp = 0
	case reflect.Int:
		intTmp = val.(int)
	case reflect.Int8:
		intTmp = int(val.(int8))
	case reflect.Int16:
		intTmp = int(val.(int16))
	case reflect.Int32:
		intTmp = int(val.(int32))
	case reflect.Int64:
		intTmp = int(val.(int64))
	case reflect.Uint:
		intTmp = int(val.(uint))
	case reflect.Uint8:
		intTmp = int(val.(uint8))
	case reflect.Uint16:
		intTmp = int(val.(uint16))
	case reflect.Uint32:
		intTmp = int(val.(uint32))
	case reflect.Float32:
		intTmp = int(val.(float32))
	case reflect.Float64:
		intTmp = int(val.(float64))
	case reflect.String:
		floatTmp, err01 := strconv.ParseFloat(val.(string), 64)
		if err01 != nil {
			intTmp = 0
		} else {
			intTmp = int(floatTmp)
		}
	default:
		intTmp = 0
	}
	return intTmp
}

// SortToQS 根据字典序排列key-value并生成url查询参数字符串
func SortToQS(data map[string]string) string {
	var (
		str  string
		keys []string
	)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if data[key] != "" {
			str += fmt.Sprintf("%s=%s&", key, data[key])
		}
	}
	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
}

// StrLen 计算字符串长度
func StrLen(str string) int {
	return len([]rune(str))
}

// ConvertToStruct 通过json序列化反序列化的方式转换不同结构体
func ConvertToStruct(in, out interface{}) (err error) {
	retByte, err := jsoniter.Marshal(in)
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(retByte, out)
	if err != nil {
		return
	}
	return
}
