package gohack

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

/*
Round 保留小数点后几位
@num: 原始数字
@precise: 在小数点后要保留的位数
*/
func Round(num float64, precise int) (output float64) {
	preciseStr := fmt.Sprintf("%d", precise)
	tpl := "%0." + preciseStr + "f"

	s := fmt.Sprintf(tpl, num)
	bitSize := 64
	output, _ = strconv.ParseFloat(s, bitSize)
	return
}

/*
InArray 判断某个值是否在数组中
@obj: 待查找的值
@target: 目标数组
*/
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

/*
StrVal 把其他类型转为string
@val: 需要转换的值
*/
func StrVal(val interface{}) (output string) {
	valueType := reflect.TypeOf(val).Kind()
	valueValue := reflect.ValueOf(val)
	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		output = fmt.Sprintf("%d", valueValue)
	case reflect.Float32, reflect.Float64:
		output = fmt.Sprintf("%f", valueValue)
	case reflect.String:
		output = fmt.Sprintf("%s", valueValue)
	}
	return
}

/*
IntVal 把其他类型转为int
@val: 需要转换的值
*/
func IntVal(val interface{}) (output int) {
	valueType := reflect.TypeOf(val).Kind()
	switch valueType {
	case reflect.Int:
		output = val.(int)
	case reflect.Int8:
		output = int(val.(int8))
	case reflect.Int16:
		output = int(val.(int16))
	case reflect.Int32:
		output = int(val.(int32))
	case reflect.Int64:
		output = int(val.(int64))
	case reflect.Uint:
		output = int(val.(uint))
	case reflect.Uint8:
		output = int(val.(uint8))
	case reflect.Uint16:
		output = int(val.(uint16))
	case reflect.Uint32:
		output = int(val.(uint32))
	case reflect.Float32:
		output = int(val.(float32))
	case reflect.Float64:
		output = int(val.(float64))
	case reflect.String:
		floatTmp, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			output = int(floatTmp)
		}
	}
	return
}

/*
SortToQS 根据字典序排列key-value并生成url查询参数字符串
@data: key-value 的字典
*/
func SortToQS(data map[string]string) (sortedQueryString string) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if data[key] != "" {
			sortedQueryString += fmt.Sprintf("%s=%s&", key, data[key])
		}
	}
	if len(sortedQueryString) > 0 {
		sortedQueryString = sortedQueryString[:len(sortedQueryString)-1]
	}
	return
}

/*
StrLen 计算字符串长度
@str: 指定字符串
*/
func StrLen(str string) int {
	return len([]rune(str))
}

/*
ConvertToStruct 通过json序列化反序列化的方式转换不同结构体
@in: 待转换的结构体
@out: 转换后的结构体
*/
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
