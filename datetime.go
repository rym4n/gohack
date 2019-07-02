package gohack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Loc 表示上海时区
var Loc *time.Location

// DefaultFmt 时间格式化的标准格式
const DefaultFmt = "2006-01-02 15:04:05"

func init() {
	Loc, _ = time.LoadLocation("Asia/Shanghai")
}

/**
 * format 格式化字符串
 * timestamp 时间戳
 * 例1：Y-m-d 返回2017-08-24
 * 例2：y年m月d日 返回 17年08月24日
 * 例3：H:i:s 返回 17:04:57
 */
func Date(format string, timestamps ...int64) string {
	var timestamp int64
	if len(timestamps) > 0 {
		timestamp = timestamps[0]
	} else {
		timestamp = Time()
	}

	//创建一个多对替换的规则
	replaceRule := strings.NewReplacer("Y", "2006", "y", "06", "m", "01", "d", "02", "H", "15", "i", "04", "s", "05")
	date := replaceRule.Replace(format)

	_timeModel := time.Unix(timestamp, 0).In(Loc)

	return _timeModel.Format(date)
}

// Time 获取当前系统时间戳
func Time() int64 {
	return Now().Unix()
}

// Now 返回上海时区的时间
func Now() time.Time {
	return time.Now().In(Loc)
}

// Today 返回当天0点时间对象
func Today() time.Time {
	n := Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, Loc)
}

// BeginOfMonth 当前月开始
func BeginOfMonth() time.Time {
	y, m, _ := Now().Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, Loc)
}

/**
 * 通过字符串获取时间戳
 * return int64 时间戳
 * return err 错误信息
 * 用例：newtimes,_ := getTime("-1days")
 * 例2：newtimes,_ := getTime("+12hours")
 * 例3：newtimes,_ := getTime("+1years")
 */
func Str2Time(str string) int64 {
	str = strings.ToLower(str)

	arrTmpMap := make(map[string]int64)
	arrTmpMap["seconds"] = 1
	arrTmpMap["minutes"] = 60
	arrTmpMap["hours"] = 3600
	arrTmpMap["days"] = 86400

	_timeNow := time.Now().In(Loc)

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

	//解析*years
	if strings.HasSuffix(str, "years") {
		numberYears, err := getFormatInt(str, "years")
		if err != nil {
			return 0
		}
		_timeNow = _timeNow.AddDate(int(numberYears), 0, 0)
		return _timeNow.Unix()
	}

	//解析 xxxx-xx-xx和xxxx-xx-xx xx:xx:xx格式的字符串
	strLen := len(str)
	if strLen == 10 {
		str = str + " 00:00:00"
	}
	strLen = len(str)
	if strLen == 19 {
		//不使用Parse解析是因为有8个小时的时差
		//_timeModel,_ := time.Parse("2006-01-02 15:04:05",str);
		_timeModel, _ := time.ParseInLocation("2006-01-02 15:04:05", str, Loc)
		return _timeModel.Unix()
	}
	return 0
}

//FormatTime convert time to a beauty format
func FormatTime(t time.Duration) string {
	var tmp float64
	timeList := [...]time.Duration{time.Hour, time.Minute, time.Second, time.Millisecond, time.Microsecond, time.Nanosecond}
	strList := [...]string{"h", "m", "s", "ms", "us", "ns"}

	for k, v := range timeList {

		if t >= v {
			tmp = float64(t) / float64(v)
			return fmt.Sprintf("%.2f%s", tmp, strList[k])
		}
	}

	return ""
}

//DiffDayNum 计算两个日期之间差多少天
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

	day := dayNumNow - dayNumStart
	dayNum = int(day)
	return dayNum

}

// Datetime2Ts 将时间格式字符串转换为时间戳
func Datetime2Ts(layout, timeString string) (int64, error) {
	t, err := time.Parse(layout, timeString)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func getFormatInt(str string, key string) (int64, error) {
	strNumber := strings.Trim(str, key)
	strNumber = strings.TrimSpace(strNumber)
	number, err := strconv.ParseInt(strNumber, 10, 0)
	if err != nil {
		return 0, errors.New("字符类型错误")
	} else {
		return number, nil
	}
}
