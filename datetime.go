package gohack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// defaultFmt 表示时间格式化的标准格式
const defaultFmt = "2006-01-02 15:04:05"

var loc *time.Location

func init() {
	//loc 初始化为上海时区
	loc, _ = time.LoadLocation("Asia/Shanghai")
}

/*
Date 格式化字符串
@timestamp: 时间戳
例1：Y-m-d 返回2017-08-24
例2：y年m月d日 返回 17年08月24日
例3：H:i:s 返回 17:04:57
*/
func Date(format string, timestamps ...int64) (formatDate string) {
	var timestamp int64
	if len(timestamps) > 0 {
		timestamp = timestamps[0]
	} else {
		timestamp = Time()
	}

	//创建一个多对替换的规则
	replaceRule := strings.NewReplacer("Y", "2006", "y", "06", "m", "01", "d", "02", "H", "15", "i", "04", "s", "05")
	date := replaceRule.Replace(format)

	formatDate = time.Unix(timestamp, 0).In(loc).Format(date)
	return
}

// Time 获取当前系统时间戳
func Time() int64 {
	return Now().Unix()
}

// Now 返回上海时区的时间
func Now() time.Time {
	return time.Now().In(loc)
}

// Today 返回当天0点时间对象
func Today() (zeroOclock time.Time) {
	n := Now()
	zeroOclock = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, loc)
	return
}

// BeginOfMonth 当前月开始
func BeginOfMonth() (firstDayOfMonth time.Time) {
	y, m, _ := Now().Date()
	firstDayOfMonth = time.Date(y, m, 1, 0, 0, 0, 0, loc)
	return
}

/*
Str2Time 通过字符串获取时间戳
例1：newtimes,_ := Str2Time("-1days")
例2：newtimes,_ := Str2Time("+12hours")
例3：newtimes,_ := Str2Time("+1years")
*/
func Str2Time(str string) int64 {
	str = strings.ToLower(str)

	arrTmpMap := make(map[string]int64)
	arrTmpMap["seconds"] = 1
	arrTmpMap["minutes"] = 60
	arrTmpMap["hours"] = 3600
	arrTmpMap["days"] = 86400

	_timeNow := time.Now().In(loc)

	//解析*seconds,*minutes,*hours,*days
	for key := range arrTmpMap {
		if strings.HasSuffix(str, key) {
			number, err := getFormatInt(str, key)
			if err != nil {
				return 0
			}
			number = _timeNow.Unix() + number*arrTmpMap[key]
			return number
		}
	}
	//解析*months
	if strings.HasSuffix(str, "months") {
		numberMonths, err := getFormatInt(str, "months")
		if err != nil {
			return 0
		}
		_timeNow = _timeNow.AddDate(0, int(numberMonths), 0)
		return _timeNow.Unix()
	}

	// 解析*years
	if strings.HasSuffix(str, "years") {
		numberYears, err := getFormatInt(str, "years")
		if err != nil {
			return 0
		}
		_timeNow = _timeNow.AddDate(int(numberYears), 0, 0)
		return _timeNow.Unix()
	}

	// 解析 xxxx-xx-xx 和 xxxx-xx-xx xx:xx:xx 格式的字符串
	strLen := len(str)
	if strLen == 10 {
		str = str + " 00:00:00"
	}
	strLen = len(str)
	if strLen == 19 {
		// 不使用Parse解析是因为有8个小时的时差
		// _timeModel,_ := time.Parse("2006-01-02 15:04:05",str);
		_timeModel, _ := time.ParseInLocation(defaultFmt, str, loc)
		return _timeModel.Unix()
	}
	return 0
}

// FormatTime 将普通时间类型转换为格式化的字符串
func FormatTime(t time.Duration) (formatTime string) {
	var tmp float64
	timeList := [...]time.Duration{time.Hour, time.Minute, time.Second, time.Millisecond, time.Microsecond, time.Nanosecond}
	strList := [...]string{"h", "m", "s", "ms", "us", "ns"}

	for k, v := range timeList {
		if t >= v {
			tmp = float64(t) / float64(v)
			formatTime = fmt.Sprintf("%.2f%s", tmp, strList[k])
			return
		}
	}

	return
}

// DiffDayNum 计算两个日期之间差多少天
func DiffDayNum(startDay string, endDay string) (dayNum int) {
	if startDay == "" || endDay == "" {
		return
	}

	daySecond := 86400
	standTime := Str2Time("1970-01-01")

	startTime := Str2Time(startDay)
	endTime := Str2Time(endDay)

	dayNumStart := float64(startTime-standTime) / float64(daySecond)
	dayNumNow := float64(endTime-standTime) / float64(daySecond)

	dayNum = int(dayNumNow - dayNumStart)
	return
}

/*
Datetime2Ts 将时间格式字符串转换为时间戳
@layout: 参考 "2006-01-02 15:04:05" 这个值做格式变换
@timeString: 时间戳字符串
*/
func Datetime2Ts(layout, timeString string) (timestamp int64, err error) {
	t, err := time.Parse(layout, timeString)
	if err == nil {
		timestamp = t.Unix()
	}
	return
}

func getFormatInt(str string, key string) (number int64, err error) {
	strNumber := strings.Trim(str, key)
	strNumber = strings.TrimSpace(strNumber)
	number, err = strconv.ParseInt(strNumber, 10, 0)
	if err != nil {
		err = errors.New("字符类型错误")
	}
	return
}
