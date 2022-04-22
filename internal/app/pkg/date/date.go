package date

import (
	"time"
)

const (
	Space      = " "
	YMD        = "2006-01-02"
	HIS        = "15:04:05"
	CommonDay  = YMD
	CommonTime = YMD + Space + HIS

	FirstSnapUp  = "06:00:00"
	SecondSnapUp = "08:30:00"
)

var Zero = time.Unix(0, 0)

// 全局设置时区
func init() {
	var local = time.FixedZone("CST", 8*3600) // 东八区
	time.Local = local
}

// Today 当天
func Today() string {
	return time.Now().Format(CommonDay)
}

func Unix(hour, min int) int64 {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local).Unix()
}

func FirstSnapUpTime() time.Time {
	str := Today() + Space + FirstSnapUp
	return ToTimeWithLayout(str, CommonTime)
}

func SecondSnapUpTime() time.Time {
	str := Today() + Space + SecondSnapUp
	return ToTimeWithLayout(str, CommonTime)
}

func ToTimeWithLayout(str, layout string) time.Time {
	if str == "" {
		return Zero
	}
	result, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		panic(err)
	}
	return result.In(time.Local)
}
