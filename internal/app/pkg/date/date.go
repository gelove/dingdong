package date

import (
	"time"
)

const (
	Local      = "Asia/Shanghai"
	Space      = " "
	YMD        = "2006-01-02"
	HIS        = "15:04:05"
	CommonDay  = YMD
	CommonTime = YMD + Space + HIS

	FirstSnapUp  = "06:00:00"
	SecondSnapUp = "08:30:00"
)

var Zero = time.Unix(0, 0)

// Today 当天
func Today() string {
	return time.Now().Format(CommonDay)
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
	loc, err := time.LoadLocation(Local)
	if err != nil {
		panic(err)
	}
	result, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		panic(err)
	}
	return result.In(loc)
}
